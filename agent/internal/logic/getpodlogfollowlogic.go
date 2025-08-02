package logic

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPodLogFollowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *commonsvc.K8sService
}

func NewGetPodLogFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPodLogFollowLogic {
	return &GetPodLogFollowLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: commonsvc.NewK8sService(ctx, svcCtx.Clientset, svcCtx.Dynamicclient, svcCtx.TektonClient, svcCtx.TriggerClient),
	}
}

func (l *GetPodLogFollowLogic) GetPodLogFollow(in *agent.PodLogRequest, stream agent.LizardAgent_GetPodLogFollowServer) (err error) {
	if _, err = l.K8sService.GetPodLog(in.Namespace, in.Podname, in.ContainerName, int64(in.Lines), true, in.Timestamps, stream, nil); err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}
