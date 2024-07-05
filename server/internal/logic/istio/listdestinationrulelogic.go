package istio

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"istio.io/client-go/pkg/apis/networking/v1beta1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDestinationRuleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListDestinationRuleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDestinationRuleLogic {
	return &ListDestinationRuleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListDestinationRuleLogic) ListDestinationRule(req *types.ListWorkloadReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	if rpcResponse, err = ag.ListDestinationRule(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.ListDestinationRule"), &agent.ListResourceRequest{
		Namespace:     req.Namespace,
		LabelSelector: req.LabelSelector,
	}); err != nil {
		l.Logger.Error(err)
		return
	}
	var r []*v1beta1.DestinationRule
	json.Unmarshal(rpcResponse.Data, &r)
	if r == nil {
		r = make([]*v1beta1.DestinationRule, 0)
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: r,
	}
	return
}
