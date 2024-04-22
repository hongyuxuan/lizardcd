package handler

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
)

func StartNacosWatch(svcCtx *svc.ServiceContext) {
	for {
		services, err := fetchAllServices(svcCtx)
		if err != nil {
			logx.Errorf("Failed to get serivces from nacos: %v", err)
			os.Exit(0)
		}
		for _, service := range services {
			if !strings.HasPrefix(service, svcCtx.Config.ServicePrefix+"lizardcd-agent") {
				continue
			}
			if _, ok := svcCtx.AgentList[service]; !ok {
				logx.Infof("A new lizardcd-agent: %s registered into nacos", service)
				cli, err := zrpc.NewClient(zrpc.RpcClientConf{
					Timeout: 5000,
					Target: fmt.Sprintf("nacos://%s:%s@%s/%s?namespaceid=%s&group=%s",
						svcCtx.Config.Nacos.Username,
						svcCtx.Config.Nacos.Password,
						svcCtx.Config.Nacos.Address,
						service,
						svcCtx.Config.Nacos.NamespaceId,
						svcCtx.Config.Nacos.Group),
				})
				if err != nil {
					logx.Error(err)
					continue
				}
				svcCtx.AgentList[service] = &types.RpcAgent{
					Client:        lizardagent.NewLizardAgent(cli),
					ServiceSource: "nacos",
				}
			}
		}
		// remove agent if services deregistered
		agentKeys := lo.Keys(svcCtx.AgentList)
		diff, _ := lo.Difference(agentKeys, services)
		for _, service := range diff {
			delete(svcCtx.AgentList, service)
			logx.Infof("Lizardcd-agent: %s removed from etcd", service)
		}
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func fetchAllServices(svcCtx *svc.ServiceContext) (services []string, err error) {
	var pageNo uint32 = 1
	var pageSize uint32 = 3
	for {
		var res model.ServiceList
		if res, err = svcCtx.NacosClient.GetAllServicesInfo(vo.GetAllServiceInfoParam{
			NameSpace: svcCtx.Config.Nacos.NamespaceId,
			GroupName: svcCtx.Config.Nacos.Group,
			PageNo:    pageNo,
			PageSize:  pageSize,
		}); err != nil {
			return
		}
		if len(res.Doms) == 0 {
			break
		}
		services = append(services, res.Doms...)
		pageNo += 1
	}
	return
}
