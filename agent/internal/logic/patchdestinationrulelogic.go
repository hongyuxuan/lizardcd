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

type PatchDestinationRuleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	istioService *svc.IstioService
}

func NewPatchDestinationRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchDestinationRuleLogic {
	return &PatchDestinationRuleLogic{
		ctx:          ctx,
		svcCtx:       svcCtx,
		Logger:       logx.WithContext(ctx),
		istioService: svc.GetIstioService(ctx, svcCtx),
	}
}

func (l *PatchDestinationRuleLogic) PatchDestinationRule(in *agent.IstioPatchRequest) (*agent.Response, error) {
	res, err := l.istioService.PatchDestinationRule(in.Namespace, in.Name, in.ResourceBody)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data, _ := json.Marshal(res)
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
