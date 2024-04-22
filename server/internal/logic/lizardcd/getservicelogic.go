package lizardcd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type GetserviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetserviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetserviceLogic {
	return &GetserviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetserviceLogic) Getservice(req *types.GetServiceReq) (resp *types.Response, err error) {
	var res *clientv3.GetResponse
	if res, err = l.svcCtx.EtcdClient.Get(l.ctx, req.ServiceName, clientv3.WithPrefix()); err != nil {
		logx.Error(err)
		return
	}
	var keymaps []map[string]interface{}
	for _, kv := range res.Kvs {
		key := utils.GetLizardAgentKey(kv.Key)
		keymaps = append(keymaps, map[string]interface{}{
			"ServiceID":   fmt.Sprintf("%s-%s", key, string(kv.Value)),
			"ServiceName": key,
			"ServiceMeta": utils.GetServiceMata(key),
		})
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: keymaps,
	}
	return
}
