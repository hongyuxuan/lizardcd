package tekton

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/samber/lo"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	tektonv1beta1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1beta1"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListTektonResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListTektonResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListTektonResourceLogic {
	return &ListTektonResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListTektonResourceLogic) ListTektonResource(req *types.ListResourceReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	if rpcResponse, err = ag.ListTektonResource(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.ListTektonResource"), &agent.TektonListRequest{
		Namespace:     req.Namespace,
		ResourceType:  req.ResourceType,
		LabelSelector: req.LabelSelector,
	}); err != nil {
		l.Logger.Error(err)
		return
	}
	switch req.ResourceType {
	case "tasks":
		r := []tektonv1.Task{}
		json.Unmarshal(rpcResponse.Data, &r)
		r = lo.Map(r, func(item tektonv1.Task, _ int) tektonv1.Task {
			item.Spec = tektonv1.TaskSpec{}
			return item
		})
		resp = &types.Response{
			Code: http.StatusOK,
			Data: r,
		}
	case "pipelines":
		r := []tektonv1.Pipeline{}
		json.Unmarshal(rpcResponse.Data, &r)
		r = lo.Map(r, func(item tektonv1.Pipeline, _ int) tektonv1.Pipeline {
			item.Spec = tektonv1.PipelineSpec{}
			return item
		})
		resp = &types.Response{
			Code: http.StatusOK,
			Data: r,
		}
	case "pipelineruns":
		r := []tektonv1.PipelineRun{}
		json.Unmarshal(rpcResponse.Data, &r)
		resp = &types.Response{
			Code: http.StatusOK,
			Data: r,
		}
	case "triggerbindings":
		r := []tektonv1beta1.TriggerBinding{}
		json.Unmarshal(rpcResponse.Data, &r)
		r = lo.Map(r, func(item tektonv1beta1.TriggerBinding, _ int) tektonv1beta1.TriggerBinding {
			item.Spec = tektonv1beta1.TriggerBindingSpec{}
			return item
		})
		resp = &types.Response{
			Code: http.StatusOK,
			Data: r,
		}
	case "triggertemplates":
		r := []tektonv1beta1.TriggerTemplate{}
		json.Unmarshal(rpcResponse.Data, &r)
		r = lo.Map(r, func(item tektonv1beta1.TriggerTemplate, _ int) tektonv1beta1.TriggerTemplate {
			item.Spec = tektonv1beta1.TriggerTemplateSpec{}
			return item
		})
		resp = &types.Response{
			Code: http.StatusOK,
			Data: r,
		}
	case "eventlisteners":
		r := []tektonv1beta1.EventListener{}
		json.Unmarshal(rpcResponse.Data, &r)
		r = lo.Map(r, func(item tektonv1beta1.EventListener, _ int) tektonv1beta1.EventListener {
			item.Spec = tektonv1beta1.EventListenerSpec{}
			return item
		})
		resp = &types.Response{
			Code: http.StatusOK,
			Data: r,
		}
	}
	return
}
