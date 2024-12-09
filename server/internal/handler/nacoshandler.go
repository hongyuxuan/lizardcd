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
	"github.com/nacos-group/nacos-sdk-go/util"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

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
				cli, err := getZrpcClient(svcCtx, service)
				if err != nil {
					logx.Error(err)
					continue
				}
				logx.Infof("A new lizardcd-agent: %s registered into nacos", service)
				svcCtx.AgentList[service] = &types.RpcAgent{
					Client:        lizardagent.NewLizardAgent(cli),
					ServiceSource: "nacos",
					Cli:           cli,
				}
				// subscribe servcie for changes
				if err = svcCtx.NacosClient.Subscribe(&vo.SubscribeParam{
					ServiceName: service,
					GroupName:   svcCtx.Config.Nacos.Group,
					SubscribeCallback: func(svcs []model.SubscribeService, e error) {
						logx.Infof("Nacos service \"%s\" changed to: %+v", service, util.ToJsonString(svcs))
						if cli, err = getZrpcClient(svcCtx, service); err != nil {
							logx.Error(err)
						} else {
							svcCtx.AgentList[service] = &types.RpcAgent{
								Client:        lizardagent.NewLizardAgent(cli),
								ServiceSource: "nacos",
								Cli:           cli,
							}
						}
					},
				}); err != nil {
					logx.Errorf("Subscribe service \"%s\" failed: %v", service, err)
				}
			}
		}
		// remove agent if services deregistered
		for k, v := range svcCtx.AgentList {
			_, ok := lo.Find(services, func(s string) bool {
				return k == s
			})
			if v.ServiceSource == "nacos" && !ok {
				delete(svcCtx.AgentList, k)
				svcCtx.NacosClient.Unsubscribe(&vo.SubscribeParam{
					ServiceName:       k,
					GroupName:         svcCtx.Config.Nacos.Group,
					SubscribeCallback: func(svcs []model.SubscribeService, e error) {},
				})
				logx.Infof("Lizardcd-agent: %s removed from nacos and unsubscribed", k)
			}
		}
		time.Sleep(time.Duration(10) * time.Second)
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

func getZrpcClient(svcCtx *svc.ServiceContext, service string) (cli zrpc.Client, err error) {
	clientConf := zrpc.RpcClientConf{
		Timeout: svcCtx.Config.Rpc.Timeout,
		Target: fmt.Sprintf("nacos://%s:%s@%s/%s?namespaceid=%s&group=%s",
			svcCtx.Config.Nacos.Username,
			svcCtx.Config.Nacos.Password,
			svcCtx.Config.Nacos.Address,
			service,
			svcCtx.Config.Nacos.NamespaceId,
			svcCtx.Config.Nacos.Group),
	}
	return zrpc.NewClient(clientConf, zrpc.WithDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                time.Duration(svcCtx.Config.Rpc.KeepaliveTime) * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: true,
	})))
}
