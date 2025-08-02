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

type GetyamlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *commonsvc.K8sService
}

func NewGetYamlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetyamlLogic {
	return &GetyamlLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: commonsvc.NewK8sService(ctx, svcCtx.Clientset, svcCtx.Dynamicclient, svcCtx.TektonClient, svcCtx.TriggerClient),
	}
}

func (l *GetyamlLogic) GetYaml(in *agent.GetYamlRequest) (resp *agent.YamlResponse, err error) {
	var res string
	if res, err = l.K8sService.GetResourceYAML(in.Namespace, in.ResourceType, in.ResourceName); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp = &agent.YamlResponse{
		Code: uint32(codes.OK),
		Data: res,
	}
	return
}
