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

type RolloutWorkloadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *commonsvc.K8sService
}

func NewRolloutWorkloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RolloutWorkloadLogic {
	return &RolloutWorkloadLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: commonsvc.NewK8sService(ctx, svcCtx.Clientset, svcCtx.Dynamicclient, svcCtx.TektonClient, svcCtx.TriggerClient),
	}
}

func (l *RolloutWorkloadLogic) RolloutWorkload(in *agent.PatchWorkloadRequest) (*agent.Response, error) {
	l.Logger.Infof("Rollout restart namespace=%s %s=%s", in.Namespace, in.WorkloadType, in.WorkloadName)
	data, err := l.K8sService.RolloutWorkload(in.Namespace, in.WorkloadType, in.WorkloadName)
	if err != nil {
		l.Logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
