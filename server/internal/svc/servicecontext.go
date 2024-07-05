package svc

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	capi "github.com/hashicorp/consul/api"
	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/config"
	"github.com/hongyuxuan/lizardcd/server/internal/middleware"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	AgentList    map[string]*types.RpcAgent
	EtcdClient   *clientv3.Client
	ConsulClient *capi.Client
	NacosClient  naming_client.INamingClient
	Sqlite       *gorm.DB
	Version      string
	Validateuser rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	svcCtx := &ServiceContext{
		Config:       c,
		AgentList:    make(map[string]*types.RpcAgent),
		Sqlite:       utils.NewSQLite(c.Sqlite, c.Log.Level),
		Validateuser: middleware.NewValidateuserMiddleware().Handle,
	}
	if c.Etcd.Address != "" {
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
		svcCtx.EtcdClient = client
	}
	if c.Consul.Address != "" {
		svcCtx.ConsulClient = utils.CreateConsul(c.Consul.Address)
	}
	if c.Nacos.Address != "" {
		fields := strings.Split(c.Nacos.Address, ":")
		host := fields[0]
		port, _ := strconv.ParseUint(fields[1], 10, 64)
		client, err := clients.NewNamingClient(
			vo.NacosClientParam{
				ServerConfigs: []constant.ServerConfig{
					{
						IpAddr: host,
						Port:   port,
					},
				},
				ClientConfig: &constant.ClientConfig{
					NamespaceId:         c.Nacos.NamespaceId,
					NotLoadCacheAtStart: true,
					LogDir:              "./logs",
					Username:            c.Nacos.Username,
					Password:            c.Nacos.Password,
				},
			},
		)
		if err != nil {
			logx.Errorf("Failed to connect to nacos: %v", err)
			os.Exit(0)
		}
		logx.Infof("Connect to nacos host=%s success", c.Nacos.Address)
		svcCtx.NacosClient = client
	}
	return svcCtx
}

func (s *ServiceContext) GetAgent(cluster, namespace string) (agent lizardagent.LizardAgent, err error) {
	for k, v := range s.AgentList {
		re, _ := regexp.Compile(k)
		if re.MatchString(fmt.Sprintf("%slizardcd-agent.%s.%s", s.Config.ServicePrefix, namespace, cluster)) {
			return v.Client, nil
		}
	}
	return nil, errorx.NewDefaultError("Cannot find lizardcd-agent of cluster=%s namespace=%s, maybe the server cannot communicated with the agent", cluster, namespace)
}

func (s *ServiceContext) GetTargetAgent(ip string) (agent lizardagent.LizardAgent, err error) {
	for k, v := range s.AgentList {
		re, _ := regexp.Compile(fmt.Sprintf("%slizardcd-agent_vm\\.(.+?)\\.%s", s.Config.ServicePrefix, ip))
		if re.MatchString(k) {
			return v.Client, nil
		}
	}
	return nil, errorx.NewDefaultError("Cannot find lizardcd-agent of ip=%s, maybe the server cannot communicated with the agent", ip)
}

func (s *ServiceContext) SetVersion(version string) {
	s.Version = version
}

func (s *ServiceContext) GetHelmSettings(tenant string) (wait bool, timeout int64, err error) {
	var settings []commontypes.Settings
	if err = s.Sqlite.Where("tenant = ?", tenant).Find(&settings).Error; err != nil {
		return
	}
	for _, s := range settings {
		if s.SettingKey == "helm_wait" {
			if s.SettingValue == "true" {
				wait = true
			} else {
				wait = false
			}
		}
		if s.SettingKey == "helm_timeout" {
			timeout, _ = strconv.ParseInt(s.SettingValue, 10, 64)
		}
	}
	return
}
