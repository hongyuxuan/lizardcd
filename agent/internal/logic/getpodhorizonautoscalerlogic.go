package logic

import (
	"context"
	"encoding/json"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPodHorizonAutoscalerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *commonsvc.K8sService
}

func NewGetPodHorizonAutoscalerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPodHorizonAutoscalerLogic {
	return &GetPodHorizonAutoscalerLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: commonsvc.NewK8sService(ctx, svcCtx.Clientset, svcCtx.Dynamicclient, svcCtx.TektonClient, svcCtx.TriggerClient),
	}
}

func (l *GetPodHorizonAutoscalerLogic) GetPodHorizonAutoscaler(in *agent.PatchWorkloadRequest) (*agent.Response, error) {
	res, err := l.K8sService.GetPodHorizonAutoscaler(in.Namespace, in.WorkloadName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data, _ := json.Marshal(res)
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
