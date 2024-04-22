package kubernetes

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScaleStatefulsetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScaleStatefulsetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScaleStatefulsetLogic {
	return &ScaleStatefulsetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScaleStatefulsetLogic) ScaleStatefulset(req *types.ScaleReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	for _, workload := range req.Workloads {
		if workload.Disabled == true {
			continue
		}
		if _, err = ag.ScaleStatefulset(context.WithValue(l.ctx, "SpanName", "rpc.ScaleStatefulset"), &agent.ScaleRequest{
			Namespace:    req.Namespace,
			WorkloadName: workload.Name,
			Replicas:     uint32(workload.Replicas),
		}); err != nil {
			l.Logger.Error(err)
			return
		}
	}
	resp = &types.Response{
		Code: http.StatusOK,
	}
	return
}
