package kubernetes

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	autov2 "k8s.io/api/autoscaling/v2"

	"github.com/zeromicro/go-zero/core/logx"
)

type GethpaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGethpaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GethpaLogic {
	return &GethpaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GethpaLogic) Gethpa(req *types.RolloutReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	if rpcResponse, err = ag.GetPodHorizonAutoscaler(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetPodHorizonAutoscaler"), &agent.RolloutWorkloadRequest{
		Namespace:    req.Namespace,
		WorkloadName: req.WorkloadName,
	}); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return &types.Response{
				Code: http.StatusOK,
			}, nil
		} else {
			l.Logger.Error(err)
			return
		}
	}
	var r *autov2.HorizontalPodAutoscaler
	json.Unmarshal(rpcResponse.Data, &r)
	resp = &types.Response{
		Code: http.StatusOK,
		Data: r,
	}
	return
}
