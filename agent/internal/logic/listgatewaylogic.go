package logic

import (
	"context"
	"encoding/json"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListGatewayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	istioService *svc.IstioService
}

func NewListGatewayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGatewayLogic {
	return &ListGatewayLogic{
		ctx:          ctx,
		svcCtx:       svcCtx,
		Logger:       logx.WithContext(ctx),
		istioService: svc.GetIstioService(ctx, svcCtx),
	}
}

func (l *ListGatewayLogic) ListGateway(in *agent.ListResourceRequest) (*agent.Response, error) {
	res, err := l.istioService.ListGateway(in.Namespace, in.LabelSelector)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data, _ := json.Marshal(res)
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
