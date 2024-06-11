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

type CreateVirtualServiceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	istioService *svc.IstioService
}

func NewCreateVirtualServiceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateVirtualServiceLogic {
	return &CreateVirtualServiceLogic{
		ctx:          ctx,
		svcCtx:       svcCtx,
		Logger:       logx.WithContext(ctx),
		istioService: svc.GetIstioService(ctx, svcCtx),
	}
}

func (l *CreateVirtualServiceLogic) CreateVirtualService(in *agent.IstioCreateRequest) (*agent.Response, error) {
	var virtualservice *v1beta1.VirtualService
	if err := json.Unmarshal(in.ResourceBody, &virtualservice); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	res, err := l.istioService.CreateVirtualService(in.Namespace, virtualservice)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data, _ := json.Marshal(res)
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
