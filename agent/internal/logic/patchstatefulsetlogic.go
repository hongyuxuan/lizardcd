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

type PatchStatefulsetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *svc.K8sService
}

func NewPatchStatefulsetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchStatefulsetLogic {
	return &PatchStatefulsetLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: svc.GetK8sService(ctx, svcCtx),
	}
}

func (l *PatchStatefulsetLogic) PatchStatefulset(in *agent.PatchWorkloadRequest) (*agent.Response, error) {
	l.Logger.Infof("patch namespace=%s statefulset[%v] container=%v image=%v", in.Namespace, in.WorkloadName, in.Container, in.Image)
	res, err := l.K8sService.PatchStatefulset(in.Namespace, in.WorkloadName, in.Container, in.Image)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data, _ := json.Marshal(res)
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
