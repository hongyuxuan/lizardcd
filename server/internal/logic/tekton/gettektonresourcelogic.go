package tekton

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	tektonv1beta1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1beta1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTektonResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTektonResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTektonResourceLogic {
	return &GetTektonResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTektonResourceLogic) GetTektonResource(req *types.ResourceReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	if rpcResponse, err = ag.GetTektonResource(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.ListTektonResource"), &agent.TektonYamlRequest{
		Namespace:    req.Namespace,
		ResourceType: req.ResourceType,
		ResourceName: req.ResourceName,
	}); err != nil {
		l.Logger.Error(err)
		return &types.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}
	switch req.ResourceType {
	case "tasks":
		var r tektonv1.Task
		json.Unmarshal(rpcResponse.Data, &r)
		r.ManagedFields = nil
		resp = &types.Response{
			Code: http.StatusOK,
			Data: r,
		}
	case "pipelines":
		var r tektonv1.Pipeline
		json.Unmarshal(rpcResponse.Data, &r)
		r.ManagedFields = nil
		resp = &types.Response{
			Code: http.StatusOK,
			Data: r,
		}
	case "pipelineruns":
		var r tektonv1.PipelineRun
		json.Unmarshal(rpcResponse.Data, &r)
		r.ManagedFields = nil
		r.Status = tektonv1.PipelineRunStatus{}
		resp = &types.Response{
			Code: http.StatusOK,
			Data: r,
		}
	case "triggerbindings":
		var r tektonv1beta1.TriggerBinding
		json.Unmarshal(rpcResponse.Data, &r)
		r.ManagedFields = nil
		resp = &types.Response{
			Code: http.StatusOK,
			Data: r,
		}
	case "triggertemplates":
		var r tektonv1beta1.TriggerTemplate
		json.Unmarshal(rpcResponse.Data, &r)
		r.ManagedFields = nil
		resp = &types.Response{
			Code: http.StatusOK,
			Data: r,
		}
	case "eventlisteners":
		var r tektonv1beta1.EventListener
		json.Unmarshal(rpcResponse.Data, &r)
		r.ManagedFields = nil
		resp = &types.Response{
			Code: http.StatusOK,
			Data: r,
		}
	}
	return
}
