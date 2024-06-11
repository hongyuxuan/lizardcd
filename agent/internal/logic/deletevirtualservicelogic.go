package logic

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteVirtualServiceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	istioService *svc.IstioService
}

func NewDeleteVirtualServiceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteVirtualServiceLogic {
	return &DeleteVirtualServiceLogic{
		ctx:          ctx,
		svcCtx:       svcCtx,
		Logger:       logx.WithContext(ctx),
		istioService: svc.GetIstioService(ctx, svcCtx),
	}
}

func (l *DeleteVirtualServiceLogic) DeleteVirtualService(in *agent.IstioGetRequest) (*agent.Response, error) {
	if err := l.istioService.DeleteVirtualService(in.Namespace, in.Name); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &agent.Response{
		Code: uint32(codes.OK),
	}, nil
}
