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

type PatchTektonResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *commonsvc.K8sService
}

func NewPatchTektonResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchTektonResourceLogic {
	return &PatchTektonResourceLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: commonsvc.NewK8sService(ctx, svcCtx.Clientset, svcCtx.Dynamicclient, svcCtx.TektonClient, svcCtx.TriggerClient),
	}
}

func (l *PatchTektonResourceLogic) PatchTektonResource(in *agent.TektonPatchRequest) (*agent.Response, error) {
	data, err := l.K8sService.TektonService.PatchResource(in.Namespace, in.ResourceType, in.ResourceName, in.PatchOptions)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
