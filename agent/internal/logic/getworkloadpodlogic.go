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

type GetWorkloadPodLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *commonsvc.K8sService
}

func NewGetWorkloadPodLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkloadPodLogic {
	return &GetWorkloadPodLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: commonsvc.NewK8sService(ctx, svcCtx.Clientset, svcCtx.Dynamicclient, svcCtx.TektonClient, svcCtx.TriggerClient),
	}
}

func (l *GetWorkloadPodLogic) GetWorkloadPod(in *agent.GetResourceRequest) (*agent.Response, error) {
	_, labels, err := l.K8sService.GetResource(in.Namespace, in.ResourceType, in.ResourceName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	res, _, err := l.K8sService.ListPods(in.Namespace, labels, "", "", 0)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data, _ := json.Marshal(res)
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
