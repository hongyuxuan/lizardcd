package lizardcd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type GetserviceLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	consulUtil *utils.ConsulUtil
}

func NewGetserviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetserviceLogic {
	return &GetserviceLogic{
		Logger:     logx.WithContext(ctx),
		ctx:        ctx,
		svcCtx:     svcCtx,
		consulUtil: utils.NewConsulUtil(ctx, svcCtx.ConsulClient),
	}
}

func (l *GetserviceLogic) Getservice(req *types.GetServiceReq) (resp *types.Response, err error) {
	agent := l.svcCtx.AgentList[req.ServiceName]
	if agent.ServiceSource == "etcd" {
		var res *clientv3.GetResponse
		if res, err = l.svcCtx.EtcdClient.Get(l.ctx, req.ServiceName, clientv3.WithPrefix()); err != nil {
			logx.Error(err)
			return
		}
		var keymaps []map[string]interface{}
		for _, kv := range res.Kvs {
			key := utils.GetLizardAgentKey(kv.Key)
			meta, _ := utils.GetServiceMata(l.svcCtx.Config.ServicePrefix, key)
			keymaps = append(keymaps, map[string]interface{}{
				"ServiceID":   fmt.Sprintf("%s-%s", key, string(kv.Value)),
				"ServiceName": key,
				"ServiceMeta": meta,
			})
		}
		resp = &types.Response{
			Code: http.StatusOK,
			Data: keymaps,
		}
	}
	if agent.ServiceSource == "consul" {
		service, err := l.consulUtil.GetService(req.ServiceName)
		if err != nil {
			l.Logger.Error(err)
			return nil, err
		}
		resp = &types.Response{
			Code: http.StatusOK,
			Data: service,
		}
	}
	if agent.ServiceSource == "nacos" {
		var res model.Service
		res, err = l.svcCtx.NacosClient.GetService(vo.GetServiceParam{
			ServiceName: req.ServiceName,
			GroupName:   l.svcCtx.Config.Nacos.Group,
		})
		var services []map[string]interface{}
		for _, host := range res.Hosts {
			services = append(services, map[string]interface{}{
				"ServiceID":   fmt.Sprintf("%s-%s:%d", req.ServiceName, host.Ip, host.Port),
				"ServiceName": req.ServiceName,
				"ServiceMeta": host.Metadata,
			})
		}
		resp = &types.Response{
			Code: http.StatusOK,
			Data: services,
		}
	}
	return
}
