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

type GetTektonResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTektonResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTektonResourceLogic {
	return &GetTektonResourceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTektonResourceLogic) GetTektonResource(in *agent.TektonYamlRequest) (*agent.Response, error) {
	var data []byte
	switch in.ResourceType {
	case "tasks":
		res, err := l.svcCtx.TektonClient.Task(in.Namespace).Get(l.ctx, in.ResourceName)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	case "pipelines":
		res, err := l.svcCtx.TektonClient.Pipeline(in.Namespace).Get(l.ctx, in.ResourceName)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	case "pipelineruns":
		res, err := l.svcCtx.TektonClient.PipelineRun(in.Namespace).Get(l.ctx, in.ResourceName)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	case "triggerbindings":
		res, err := l.svcCtx.TektonClient.TriggerBinding(in.Namespace).Get(l.ctx, in.ResourceName)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	case "triggertemplates":
		res, err := l.svcCtx.TektonClient.TriggerTemplate(in.Namespace).Get(l.ctx, in.ResourceName)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	case "eventlisteners":
		res, err := l.svcCtx.TektonClient.EventListener(in.Namespace).Get(l.ctx, in.ResourceName)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	}
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
