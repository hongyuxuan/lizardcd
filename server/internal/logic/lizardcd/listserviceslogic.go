package lizardcd

import (
	"context"
	"net/http"
	"strings"

	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/samber/lo"

	"github.com/zeromicro/go-zero/core/logx"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ListservicesLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	consulUtil *utils.ConsulUtil
}

func NewListservicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListservicesLogic {
	return &ListservicesLogic{
		Logger:     logx.WithContext(ctx),
		ctx:        ctx,
		svcCtx:     svcCtx,
		consulUtil: utils.NewConsulUtil(ctx, svcCtx.ConsulClient),
	}
}

func (l *ListservicesLogic) Listservices() (resp *types.Response, err error) {
	var services []map[string]string
	if l.svcCtx.Config.Consul.Address != "" {
		var res map[string][]string
		if res, err = l.consulUtil.ListServices(); err != nil {
			l.Logger.Error(err)
			return
		}
		for k := range res {
			if strings.HasPrefix(k, "lizardcd-agent") {
				services = append(services, map[string]string{
					"service_name": k,
					"service_type": "consul",
				})
			}
		}
	}
	if l.svcCtx.Config.Etcd.Address != "" {
		var res *clientv3.GetResponse
		if res, err = l.svcCtx.EtcdClient.Get(l.ctx, "lizardcd-agent", clientv3.WithPrefix()); err != nil {
			logx.Error(err)
			return
		}
		var serviceList []string
		for _, kv := range res.Kvs {
			key := utils.GetLizardAgentKey(kv.Key)
			serviceList = append(serviceList, key)
		}
		serviceList = lo.Uniq(serviceList)
		services = lo.Map(serviceList, func(x string, _ int) map[string]string {
			return map[string]string{
				"service_name":   x,
				"service_source": "etcd",
			}
		})
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: services,
	}
	return
}
