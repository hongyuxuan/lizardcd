package logic

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *commonsvc.K8sService
}

func NewDeleteResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteResourceLogic {
	return &DeleteResourceLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: commonsvc.NewK8sService(ctx, svcCtx.Clientset, svcCtx.Dynamicclient, svcCtx.TektonClient, svcCtx.TriggerClient),
	}
}

func (l *DeleteResourceLogic) DeleteResource(in *agent.DeleteResourceRequest) (*agent.Response, error) {
	l.Logger.Infof("Delete namespace=%s %s[%v]", in.Namespace, in.ResourceType, in.ResourceName)
	if err := l.K8sService.DeleteResource(in.Namespace, in.ResourceType, in.ResourceName, in.Force); err != nil {
		return nil, err
	}
	return &agent.Response{
		Code: uint32(codes.OK),
	}, nil
}
