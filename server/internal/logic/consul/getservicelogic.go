package consul

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
	service, err := l.consulUtil.GetService(req.ServiceName)
	if err != nil {
		l.Logger.Error(err)
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: service,
	}
	return
}
