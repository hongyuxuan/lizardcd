package logic

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDeploymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *svc.K8sService
}

func NewDeleteDeploymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDeploymentLogic {
	return &DeleteDeploymentLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: svc.GetK8sService(ctx, svcCtx),
	}
}

func (l *DeleteDeploymentLogic) DeleteDeployment(in *agent.GetWorkloadRequest) (*agent.Response, error) {
	l.Logger.Infof("Delete namespace=%s deployment[%v]", in.Namespace, in.WorkloadName)
	if err := l.K8sService.DeleteDeployment(in.Namespace, in.WorkloadName); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &agent.Response{
		Code: uint32(codes.OK),
	}, nil
}
