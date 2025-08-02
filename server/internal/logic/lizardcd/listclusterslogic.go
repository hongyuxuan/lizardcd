package lizardcd

import (
	"context"
	"net/http"
	"strings"

	"github.com/hongyuxuan/lizardcd/common/constant"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/samber/lo"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListclustersLogic struct {
	logx.Logger
	ctx               context.Context
	svcCtx            *svc.ServiceContext
	listservicesLogic *ListservicesLogic
}

func NewListclustersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListclustersLogic {
	return &ListclustersLogic{
		Logger:            logx.WithContext(ctx),
		ctx:               ctx,
		svcCtx:            svcCtx,
		listservicesLogic: NewListservicesLogic(ctx, svcCtx),
	}
}

func (l *ListclustersLogic) Listclusters(req *types.ListClusterReq) (resp *types.Response, err error) {
	_, role, _, namespaces := utils.GetPayload(l.ctx)
	var clusterMap = make(map[string][]string)
	for _, svc := range l.svcCtx.ListServices(role, namespaces) {
		l.Logger.Debug(svc)
		if strings.Contains(svc.Service, "agent_vm") {
			if !req.WithVm {
				continue
			} else {
				svc.Cluster = "vm"
			}
		}
		if _, ok := clusterMap[svc.Cluster]; !ok {
			clusterMap[svc.Cluster] = []string{}
		}
		if svc.Namespace == "*" {
			nss := l.svcCtx.GetNamespaces(svc.ServiceName)
			if role == constant.ROLE_ADMIN {
				clusterMap[svc.Cluster] = append(clusterMap[svc.Cluster], nss...)
			} else {
				if found := lo.Intersect(nss, namespaces); len(found) > 0 {
					clusterMap[svc.Cluster] = append(clusterMap[svc.Cluster], found...)
				}
			}
		} else {
			if _, ok := lo.Find(namespaces, func(s string) bool {
				return s == svc.Namespace
			}); ok || role == constant.ROLE_ADMIN {
				clusterMap[svc.Cluster] = append(clusterMap[svc.Cluster], svc.Namespace)
			}
		}
	}
	for k, v := range clusterMap {
		if len(v) == 0 {
			delete(clusterMap, k)
		}
		clusterMap[k] = lo.Uniq(v)
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: clusterMap,
	}
	return
}
