package svc

import (
	"fmt"
	"regexp"

	capi "github.com/hashicorp/consul/api"
	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/config"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	AgentList    map[string]lizardagent.LizardAgent
	ConsulClient *capi.Client
	Sqlite       *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		AgentList:    make(map[string]lizardagent.LizardAgent),
		ConsulClient: utils.CreateConsul(c.Consul.Address),
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
