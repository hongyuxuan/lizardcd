package consul

import (
	"context"
	"net/http"
	"strings"

	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	res, err := l.consulUtil.ListServices()
	if err != nil {
		l.Logger.Error(err)
		return
	}
	var services []map[string]string
	for k := range res {
		if strings.HasPrefix(k, "lizardcd-agent") {
			services = append(services, map[string]string{
				"service_name": k,
			})
		}
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: services,
	}
	return
}
