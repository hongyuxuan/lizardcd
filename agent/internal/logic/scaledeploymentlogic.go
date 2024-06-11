package logic

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScaleDeploymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *svc.K8sService
}

func NewScaleDeploymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScaleDeploymentLogic {
	return &ScaleDeploymentLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: svc.GetK8sService(ctx, svcCtx),
	}
}

func (l *ScaleDeploymentLogic) ScaleDeployment(in *agent.ScaleRequest) (resp *agent.Response, err error) {
	if err = l.K8sService.ScaleDeployment(in.Namespace, in.WorkloadName, in.Replicas); err != nil {
		return
	}
	resp = &agent.Response{
		Code: uint32(codes.OK),
	}
	return
}
