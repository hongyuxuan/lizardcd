package logic

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTektonEndpointLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTektonEndpointLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTektonEndpointLogic {
	return &GetTektonEndpointLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTektonEndpointLogic) GetTektonEndpoint(in *agent.TektonEndpointRequest) (*agent.Response, error) {
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: []byte(l.svcCtx.Config.TektonEndpoint),
	}, nil
}
