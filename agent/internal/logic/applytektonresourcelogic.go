package logic

import (
	"context"
	"fmt"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyTektonResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyTektonResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyTektonResourceLogic {
	return &ApplyTektonResourceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ApplyTektonResourceLogic) ApplyTektonResource(in *agent.TektonYamlRequest) (*agent.Response, error) {
	l.Logger.Debug("\n", in.Ymlstring)
	switch in.ResourceType {
	case "Task":
		if err := l.svcCtx.TektonClient.Task(in.Namespace).Create(l.ctx, in.Ymlstring); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	case "Pipeline":
		if err := l.svcCtx.TektonClient.Pipeline(in.Namespace).Create(l.ctx, in.Ymlstring); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	case "PipelineRun":
		if err := l.svcCtx.TektonClient.PipelineRun(in.Namespace).Create(l.ctx, in.Ymlstring); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	case "TriggerBinding":
		if err := l.svcCtx.TektonClient.TriggerBinding(in.Namespace).Create(l.ctx, in.Ymlstring); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	case "TriggerTemplate":
		if err := l.svcCtx.TektonClient.TriggerTemplate(in.Namespace).Create(l.ctx, in.Ymlstring); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	case "EventListener":
		if err := l.svcCtx.TektonClient.EventListener(in.Namespace).Create(l.ctx, in.Ymlstring); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	default:
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid tekton kind=%s", in.ResourceType))
	}
	return &agent.Response{
		Code: uint32(codes.OK),
	}, nil
}
