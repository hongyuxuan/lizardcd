package logic

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPodLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *svc.K8sService
}

func NewGetPodLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPodLogLogic {
	return &GetPodLogLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: svc.GetK8sService(ctx, svcCtx),
	}
}

func (l *GetPodLogLogic) GetPodLog(in *agent.PodLogRequest) (*agent.YamlResponse, error) {
	res, err := l.K8sService.GetPodLog(in.Namespace, in.Podname, in.ContainerName, int64(in.Lines), false, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &agent.YamlResponse{
		Code: uint32(codes.OK),
		Data: res,
	}, nil
}
