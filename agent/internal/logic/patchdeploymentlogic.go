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

type PatchDeploymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *svc.K8sService
}

func NewPatchDeploymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchDeploymentLogic {
	return &PatchDeploymentLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: svc.GetK8sService(ctx, svcCtx),
	}
}

func (l *PatchDeploymentLogic) PatchDeployment(in *agent.PatchWorkloadRequest) (*agent.Response, error) {
	l.Logger.Infof("Patch namespace=%s deployment[%v] container=%v image=%v", in.Namespace, in.WorkloadName, in.Container, in.Image)
	res, err := l.K8sService.PatchDeployment(in.Namespace, in.WorkloadName, in.Container, in.Image)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data, _ := json.Marshal(res)
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
