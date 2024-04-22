package svc

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	capi "github.com/hashicorp/consul/api"
	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	AgentList    map[string]lizardagent.LizardAgent
	ConsulClient *capi.Client
	EtcdClient   *clientv3.Client
	Sqlite       *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	etcdHosts := strings.Split(c.Etcd.Address, ",")
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdHosts,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logx.Errorf("Failed to connect to etcd: %v", err)
		os.Exit(0)
	}
	logx.Infof("Connect to etcd host=%s success", c.Etcd.Address)

	return &ServiceContext{
		Config:       c,
		AgentList:    make(map[string]lizardagent.LizardAgent),
		ConsulClient: utils.CreateConsul(c.Consul.Address),
		EtcdClient:   client,
		Sqlite:       utils.NewSQLite(c.Sqlite, c.Log.Level),
	}
}

func (s *ServiceContext) GetAgent(cluster, namespace string) (agent lizardagent.LizardAgent, err error) {
	for k, v := range s.AgentList {
		re, _ := regexp.Compile(k)
		if re.MatchString(fmt.Sprintf("lizardcd-agent.%s.%s", namespace, cluster)) {
			return v, nil
		}
	}
	return nil, errorx.NewDefaultError(fmt.Sprintf("Cannot find lizardcd-agent of cluster=%s namespace=%s, maybe the server cannot communicated with the agent", cluster, namespace))
}
