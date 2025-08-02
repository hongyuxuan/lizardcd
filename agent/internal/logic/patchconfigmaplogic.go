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

type PatchConfigmapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *commonsvc.K8sService
}

func NewPatchConfigmapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchConfigmapLogic {
	return &PatchConfigmapLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: commonsvc.NewK8sService(ctx, svcCtx.Clientset, svcCtx.Dynamicclient, svcCtx.TektonClient, svcCtx.TriggerClient),
	}
}

func (l *PatchConfigmapLogic) PatchConfigmap(in *agent.PatchConfigmapRequest) (*agent.Response, error) {
	data, err := l.K8sService.PatchConfig(in.Namespace, in.ResourceType, in.ResourceName, in.Key, in.Value)
	if err != nil {
		l.Logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	b, _ := json.Marshal(data)
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: b,
	}, nil
}
