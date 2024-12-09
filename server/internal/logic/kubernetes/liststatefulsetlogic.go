package kubernetes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	v1 "k8s.io/api/apps/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListStatefulsetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListStatefulsetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListStatefulsetLogic {
	return &ListStatefulsetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListStatefulsetLogic) ListStatefulset(req *types.ListWorkloadReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	if rpcResponse, err = ag.ListStatefulset(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.ListStatefulset"), &agent.ListResourceRequest{
		Namespace:     req.Namespace,
		LabelSelector: req.LabelSelector,
	}); err != nil {
		l.Logger.Error(err)
		return
	}
	var r []v1.StatefulSet
	json.Unmarshal(rpcResponse.Data, &r)
	if r == nil {
		r = make([]v1.StatefulSet, 0)
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: r,
	}
	return
}
