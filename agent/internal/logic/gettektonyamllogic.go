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

type GetTektonYamlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTektonYamlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTektonYamlLogic {
	return &GetTektonYamlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTektonYamlLogic) GetTektonYaml(in *agent.TektonYamlRequest) (*agent.YamlResponse, error) {
	var data []byte
	switch in.ResourceType {
	case "tasks":
		res, err := l.svcCtx.TektonClient.Task(in.Namespace).GetYaml(l.ctx, in.ResourceName)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	case "pipelines":
		res, err := l.svcCtx.TektonClient.Pipeline(in.Namespace).GetYaml(l.ctx, in.ResourceName)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	case "pipelineruns":
		res, err := l.svcCtx.TektonClient.PipelineRun(in.Namespace).GetYaml(l.ctx, in.ResourceName)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	case "triggerbindings":
		res, err := l.svcCtx.TektonClient.TriggerBinding(in.Namespace).GetYaml(l.ctx, in.ResourceName)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	case "triggertemplates":
		res, err := l.svcCtx.TektonClient.TriggerTemplate(in.Namespace).GetYaml(l.ctx, in.ResourceName)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	case "eventlisteners":
		res, err := l.svcCtx.TektonClient.EventListener(in.Namespace).GetYaml(l.ctx, in.ResourceName)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	}
	return &agent.YamlResponse{
		Code: uint32(codes.OK),
		Data: string(data),
	}, nil
}
