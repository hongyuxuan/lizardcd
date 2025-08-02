package lizardcd

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/constant"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
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
	services := lo.Filter(l.svcCtx.ListServices(role, namespaces), func(svc *commontypes.ServiceMeta, _ int) bool {
		if svc.Namespace != "*" {
			if _, ok := lo.Find(namespaces, func(s string) bool {
				return s == svc.Namespace
			}); !ok && role != constant.ROLE_ADMIN {
				return false
			}
			return true
		} else {
			nss := l.svcCtx.GetNamespaces(svc.ServiceName)
			if found := lo.Intersect(nss, namespaces); len(found) > 0 || role == constant.ROLE_ADMIN {
				return true
			}
			return false
		}
	})
	if len(services) == 0 {
		services = []*commontypes.ServiceMeta{}
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: services,
	}
	return
}
