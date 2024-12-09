package logic

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTektonResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteTektonResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTektonResourceLogic {
	return &DeleteTektonResourceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteTektonResourceLogic) DeleteTektonResource(in *agent.TektonYamlRequest) (*agent.Response, error) {
	switch in.ResourceType {
	case "tasks":
		if err := l.svcCtx.TektonClient.Task(in.Namespace).Delete(l.ctx, in.ResourceName); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	case "pipelines":
		if err := l.svcCtx.TektonClient.Pipeline(in.Namespace).Delete(l.ctx, in.ResourceName); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	case "pipelineruns":
		if err := l.svcCtx.TektonClient.PipelineRun(in.Namespace).Delete(l.ctx, in.ResourceName); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	case "triggerbindings":
		if err := l.svcCtx.TektonClient.TriggerBinding(in.Namespace).Delete(l.ctx, in.ResourceName); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	case "triggertemplates":
		if err := l.svcCtx.TektonClient.TriggerTemplate(in.Namespace).Delete(l.ctx, in.ResourceName); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	case "eventlisteners":
		if err := l.svcCtx.TektonClient.EventListener(in.Namespace).Delete(l.ctx, in.ResourceName); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &agent.Response{
		Code: uint32(codes.OK),
	}, nil
}
