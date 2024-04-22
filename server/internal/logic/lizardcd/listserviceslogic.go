package lizardcd

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

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
	var services []map[string]string
	for k, v := range l.svcCtx.AgentList {
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
