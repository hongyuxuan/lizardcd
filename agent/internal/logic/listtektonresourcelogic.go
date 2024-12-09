package logic

import (
	"context"
	"encoding/json"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListTektonResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListTektonResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTektonResourceLogic {
	return &ListTektonResourceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// tekton
func (l *ListTektonResourceLogic) ListTektonResource(in *agent.TektonListRequest) (*agent.Response, error) {
	var data []byte
	switch in.ResourceType {
	case "tasks":
		res, err := l.svcCtx.TektonClient.Task(in.Namespace).List(l.ctx, metav1.ListOptions{
			LabelSelector: in.LabelSelector,
		})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	case "pipelines":
		res, err := l.svcCtx.TektonClient.Pipeline(in.Namespace).List(l.ctx, metav1.ListOptions{
			LabelSelector: in.LabelSelector,
		})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	case "pipelineruns":
		res, err := l.svcCtx.TektonClient.PipelineRun(in.Namespace).List(l.ctx, metav1.ListOptions{
			LabelSelector: in.LabelSelector,
		})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	case "triggerbindings":
		res, err := l.svcCtx.TektonClient.TriggerBinding(in.Namespace).List(l.ctx, metav1.ListOptions{
			LabelSelector: in.LabelSelector,
		})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	case "triggertemplates":
		res, err := l.svcCtx.TektonClient.TriggerTemplate(in.Namespace).List(l.ctx, metav1.ListOptions{
			LabelSelector: in.LabelSelector,
		})
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		data, _ = json.Marshal(res)
	case "eventlisteners":
		res, err := l.svcCtx.TektonClient.EventListener(in.Namespace).List(l.ctx, metav1.ListOptions{
			LabelSelector: in.LabelSelector,
		})
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
