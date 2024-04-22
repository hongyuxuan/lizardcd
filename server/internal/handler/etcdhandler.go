package handler

import (
	"context"
	"os"
	"strings"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func StartEtcdWatch(svcCtx *svc.ServiceContext) {
	etcdHosts := strings.Split(svcCtx.Config.Etcd.Address, ",")

	// fetch all agents from etcd
	res, err := svcCtx.EtcdClient.Get(context.TODO(), "lizardcd-agent", clientv3.WithPrefix())
	if err != nil {
		logx.Errorf("Failed to get keys from etcd: %v", err)
		os.Exit(0)
	}
	for _, kv := range res.Kvs {
		key := utils.GetLizardAgentKey(kv.Key)
		if err := addAgentList(svcCtx, etcdHosts, key); err != nil {
			continue
		}
	}

	// start to watch etcd key
	wc := svcCtx.EtcdClient.Watch(context.TODO(), "lizardcd-agent", clientv3.WithPrefix())
	for v := range wc {
		for _, e := range v.Events {
			key := utils.GetLizardAgentKey(e.Kv.Key)
			if e.Type == mvccpb.PUT {
				if err := addAgentList(svcCtx, etcdHosts, key); err != nil {
					continue
				}
			} else if e.Type == mvccpb.DELETE {
				logx.Infof("Lizardcd-agent=%s removed from etcd", key)
				delete(svcCtx.AgentList, key)
			}
		}
	}
}

func addAgentList(svcCtx *svc.ServiceContext, etcdHosts []string, key string) error {
	if _, ok := svcCtx.AgentList[key]; !ok {
		logx.Infof("A new lizardcd-agent=%s registered into etcd", key)
		cli, err := zrpc.NewClient(zrpc.RpcClientConf{
			Etcd: discov.EtcdConf{
				Hosts: etcdHosts,
				Key:   key,
			},
		})
		if err != nil {
			logx.Error(err)
			return err
		}
		svcCtx.AgentList[key] = lizardagent.NewLizardAgent(cli)
	}
	return nil
}
