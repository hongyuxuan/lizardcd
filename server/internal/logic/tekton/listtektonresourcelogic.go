package tekton

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
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
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
	}
	var data []byte
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.ListTektonResource(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.ListTektonResource"), &agent.TektonListRequest{
			Namespace:     req.Namespace,
			ResourceType:  req.ResourceType,
			LabelSelector: req.LabelSelector,
			FieldSelector: req.FieldSelector,
			Continue:      req.Continue,
			Limit:         req.Limit,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		data = rpcResponse.Data
	} else if ks != nil && ks.IsValid() {
		if data, err = ks.TektonService.ListResource(req.Namespace, req.ResourceType, req.LabelSelector, req.FieldSelector, req.Continue, req.Limit); err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot ListTektonResource of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}
	var res interface{}
	json.Unmarshal(data, &res)
	resp.Data = res
	return
}
