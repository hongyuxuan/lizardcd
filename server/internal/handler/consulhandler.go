package handler

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	capi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type WatchHandler interface {
	Handler(uint64, interface{})
}

type ConsulWatch struct {
	watchers map[string]*watch.Plan // store plans
	RWMutex  *sync.RWMutex
	svcCtx   *svc.ServiceContext
}

func (c ConsulWatch) Handler(_ uint64, data interface{}) {
	switch d := data.(type) {
	// "services" watch type returns map[string][]string type. see:https://www.consul.io/docs/dynamic-app-config/watches#services
	case map[string][]string:
		for k := range d {
			if _, ok := c.watchers[k]; ok || k == "consul" {
				continue
			}
			if strings.HasPrefix(k, c.svcCtx.Config.ServicePrefix+"lizardcd-agent") {
				// add lizardcd-agent service to agentList
				if _, ok := c.svcCtx.AgentList[k]; !ok {
					cli, err := zrpc.NewClient(zrpc.RpcClientConf{
						Timeout: c.svcCtx.Config.Rpc.Timeout,
						Target:  fmt.Sprintf("consul://%s/%s?wait=60s", c.svcCtx.Config.Consul.Address, k),
					}, zrpc.WithDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
						Time:                time.Duration(c.svcCtx.Config.Rpc.KeepaliveTime) * time.Second,
						Timeout:             time.Second,
						PermitWithoutStream: true,
					})))
					if err != nil {
						logx.Error(err)
						continue
					}
					logx.Infof("A new lizardcd-agent: %s registered into consul", k)
					c.svcCtx.AgentList[k] = &types.RpcAgent{
						Client:        lizardagent.NewLizardAgent(cli),
						ServiceSource: "consul",
						Cli:           cli,
					}
				}
			}
			// start creating one watch plan to watch every service
			c.InsertServiceWatch(k)
		}

		// read watchers and delete deregister services
		c.RWMutex.RLock()
		defer c.RWMutex.RUnlock()
		watchers := c.watchers
		for k, plan := range watchers {
			if _, ok := d[k]; !ok {
				plan.Stop()
				delete(watchers, k)
				if strings.HasPrefix(k, c.svcCtx.Config.ServicePrefix+"lizardcd-agent") {
					delete(c.svcCtx.AgentList, k)
					logx.Infof("Lizardcd-agent: %s removed from consul", k)
				}
			}
		}
	default:
		logx.Errorf("Can't decide the watch type: %v", &d)
	}
}

func NewWatchPlan(watchType string, opts map[string]interface{}, handler WatchHandler) (*watch.Plan, error) {
	var options = map[string]interface{}{
		"type": watchType,
	}
	// combine params
	for k, v := range opts {
		options[k] = v
	}
	pl, err := watch.Parse(options)
	if err != nil {
		return nil, err
	}
	pl.Handler = handler.Handler
	return pl, nil
}

func RunWatchPlan(plan *watch.Plan, address string) error {
	defer plan.Stop()
	err := plan.Run(address)
	if err != nil {
		logx.Errorf("Run consul error: %v", err)
		return err
	}
	return nil
}

func StartConsulWatch(svcCtx *svc.ServiceContext) {
	cw := ConsulWatch{
		watchers: make(map[string]*watch.Plan),
		RWMutex:  new(sync.RWMutex),
		svcCtx:   svcCtx,
	}
	wp, err := NewWatchPlan("services", nil, cw)
	if err != nil {
		logx.Errorf("new watch plan failed: %v", err)
		os.Exit(0)
	}
	err = RunWatchPlan(wp, svcCtx.Config.Consul.Address)
	if err != nil {
		os.Exit(0)
	}
}

type ServiceWatch struct {
	Address string
}

func (s ServiceWatch) Handler(_ uint64, data interface{}) {
	switch d := data.(type) {
	case []*capi.ServiceEntry:
		for _, entry := range d {
			logx.Debugf("ServiceID=%s, instance=%s, status=%s", entry.Service.ID, entry.Service.Address, entry.Checks.AggregatedStatus())
		}
	}
}

func (c ConsulWatch) InsertServiceWatch(serviceName string) {
	serviceOpts := map[string]interface{}{
		"service": serviceName,
	}
	sw := ServiceWatch{
		Address: c.svcCtx.Config.Consul.Address,
	}
	servicePlan, err := NewWatchPlan("service", serviceOpts, sw)
	if err != nil {
		logx.Errorf("New service watch failed: %v", err)
	}

	go func() {
		_ = RunWatchPlan(servicePlan, sw.Address)
	}()
	defer c.RWMutex.Unlock()
	c.RWMutex.Lock()
	c.watchers[serviceName] = servicePlan
}
