package lizardcd

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/constant"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/samber/lo"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListservicesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListservicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListservicesLogic {
	return &ListservicesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListservicesLogic) Listservices() (resp *types.Response, err error) {
	_, role, _, namespaces := utils.GetPayload(l.ctx)
	var services []map[string]string
	for k, v := range l.svcCtx.AgentList {
		meta, _ := utils.GetServiceMata(l.svcCtx.Config.ServicePrefix, k)
		if _, ok := lo.Find(namespaces, func(s string) bool {
			if meta != nil {
				return s == meta["Namespace"]
			}
			return false
		}); !ok && role != constant.ROLE_ADMIN {
			continue
		}
		services = append(services, map[string]string{
			"service_name":   k,
			"service_source": v.ServiceSource,
		})
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: services,
	}
	return
}
