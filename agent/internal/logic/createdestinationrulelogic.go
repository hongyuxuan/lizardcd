package logic

import (
	"context"
	"encoding/json"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"istio.io/client-go/pkg/apis/networking/v1beta1"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDestinationRuleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	istioService *svc.IstioService
}

func NewCreateDestinationRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDestinationRuleLogic {
	return &CreateDestinationRuleLogic{
		ctx:          ctx,
		svcCtx:       svcCtx,
		Logger:       logx.WithContext(ctx),
		istioService: svc.GetIstioService(ctx, svcCtx),
	}
}

// istio
func (l *CreateDestinationRuleLogic) CreateDestinationRule(in *agent.IstioCreateRequest) (*agent.Response, error) {
	var destinationrule *v1beta1.DestinationRule
	if err := json.Unmarshal(in.ResourceBody, &destinationrule); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	res, err := l.istioService.CreateDestinationRule(in.Namespace, destinationrule)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data, _ := json.Marshal(res)
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
