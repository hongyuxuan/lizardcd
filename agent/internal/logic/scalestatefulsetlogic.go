package logic

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScaleStatefulsetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *svc.K8sService
}

func NewScaleStatefulsetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScaleStatefulsetLogic {
	return &ScaleStatefulsetLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: svc.GetK8sService(ctx, svcCtx),
	}
}

func (l *ScaleStatefulsetLogic) ScaleStatefulset(in *agent.ScaleRequest) (resp *agent.Response, err error) {
	if err = l.K8sService.ScaleStatefulset(in.Namespace, in.WorkloadName, in.Replicas); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp = &agent.Response{
		Code: uint32(codes.OK),
	}
	return
}
