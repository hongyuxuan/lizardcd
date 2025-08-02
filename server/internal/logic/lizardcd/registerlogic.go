package lizardcd

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterServiceReq) (resp *types.Response, err error) {
	if err = l.svcCtx.RegisterAgent(req.ServiceKey, req.Endpoint, req.Proxy, req.Kubeconfig, req.Labels); err != nil {
		l.Logger.Error(err)
		return
	}

	// 插入数据库
	if err = l.svcCtx.Database.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.SaveAgent")).Save(&commontypes.Agent{
		ServiceKey: req.ServiceKey,
		Endpoint:   req.Endpoint,
		Proxy:      req.Proxy,
		Kubeconfig: req.Kubeconfig,
		Labels:     req.Labels,
	}).Error; err != nil {
		return
	}

	return &types.Response{
		Code:    http.StatusOK,
		Message: "register success",
	}, nil
}
