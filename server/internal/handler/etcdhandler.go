package handler

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func StartEtcdWatch(svcCtx *svc.ServiceContext) {
	etcdHosts := strings.Split(svcCtx.Config.Etcd.Address, ",")

	// fetch all agents from etcd
	res, err := svcCtx.EtcdClient.Get(context.TODO(), svcCtx.Config.ServicePrefix+"lizardcd-agent", clientv3.WithPrefix())
	if err != nil {
		logx.Errorf("Failed to get keys from etcd: %v", err)
		os.Exit(0)
	}
	for _, kv := range res.Kvs {
		key := utils.GetLizardAgentKey(kv.Key)
		go addAgentList(svcCtx, etcdHosts, key)
	}

	// start to watch etcd key
	wc := svcCtx.EtcdClient.Watch(context.TODO(), svcCtx.Config.ServicePrefix+"lizardcd-agent", clientv3.WithPrefix())
	for v := range wc {
		for _, e := range v.Events {
			key := utils.GetLizardAgentKey(e.Kv.Key)
			if e.Type == mvccpb.PUT {
				go addAgentList(svcCtx, etcdHosts, key)
			} else if e.Type == mvccpb.DELETE {
				if _, ok := svcCtx.AgentList[key]; ok {
					if svcCtx.AgentList[key].ServiceSource == "etcd" {
						svcCtx.AgentList[key].Count -= 1
						if svcCtx.AgentList[key].Count == 0 {
							delete(svcCtx.AgentList, key)
							logx.Infof("Lizardcd-agent: %s removed from etcd", key)
						}
					}
				}
			}
		}
	}
}

func WatchAgentManually(svcCtx *svc.ServiceContext) {
	for {
		var agents []commontypes.Agent
		if err := svcCtx.Database.Find(&agents).Error; err != nil {
			logx.Errorf("Failed to get agents from db: %v", err)
			continue
		}
		fromDB := lo.Map(agents, func(item commontypes.Agent, _ int) string {
			return item.ServiceKey
		})
		var onlineAgent []string
		for k, v := range svcCtx.AgentList {
			if v.ServiceSource == "manual" {
				onlineAgent = append(onlineAgent, k)
			}
		}
		agentToRemove, _ := lo.Difference(onlineAgent, fromDB)
		for _, k := range agentToRemove {
			if _, ok := svcCtx.AgentList[k]; ok {
				svcCtx.AgentList[k].Cli.Conn().Close()
				delete(svcCtx.AgentList, k)
				logx.Infof("Lizardcd-agent %s removed from lizardcd-server manually", k)
			}
		}

		var onlineK8s []string
		for k := range svcCtx.K8sList {
			onlineK8s = append(onlineK8s, k)
		}
		k8sToRemove, _ := lo.Difference(onlineK8s, fromDB)
		for _, k := range k8sToRemove {
			if _, ok := svcCtx.K8sList[k]; ok {
				delete(svcCtx.K8sList, k)
				logx.Infof("K8s connection %s removed from lizardcd-server", k)
			}
		}

		for _, agent := range agents {
			if err := svcCtx.RegisterAgent(agent.ServiceKey, agent.Endpoint, agent.Proxy, agent.Kubeconfig, agent.Labels); err != nil {
				logx.Error(err)
				continue
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func addAgentList(svcCtx *svc.ServiceContext, etcdHosts []string, key string) {
	for {
		if _, ok := svcCtx.AgentList[key]; !ok {
			cli, err := zrpc.NewClient(zrpc.RpcClientConf{
				Timeout: svcCtx.Config.Rpc.Timeout,
				Etcd: discov.EtcdConf{
					Hosts: etcdHosts,
					Key:   key,
				},
			}, zrpc.WithDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{ // add keepalive option
				Time:                time.Duration(svcCtx.Config.Rpc.KeepaliveTime) * time.Second,
				Timeout:             time.Second,
				PermitWithoutStream: true,
			})))
			if err != nil {
				logx.Error(err)
				time.Sleep(time.Duration(svcCtx.Config.Rpc.RetryInterval) * time.Second) // sleep <RetryInterval> seconds and try again
				continue
			} else {
				logx.Infof("A new lizardcd-agent: %s registered into etcd", key)
				svcCtx.AgentList[key] = &types.RpcAgent{
					Client:        lizardagent.NewLizardAgent(cli),
					ServiceSource: "etcd",
					Cli:           cli,
					Count:         1,
				}
				return
			}
		}
		svcCtx.AgentList[key].Count += 1
		return
	}
}
