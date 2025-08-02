package kubernetes

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
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
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var r *autov2.HorizontalPodAutoscaler
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.GetPodHorizonAutoscaler(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetPodHorizonAutoscaler"), &agent.PatchWorkloadRequest{
			Namespace:    req.Namespace,
			WorkloadType: req.WorkloadType,
			WorkloadName: req.WorkloadName,
		}); err != nil {
			if strings.Contains(err.Error(), "could not find") || strings.Contains(err.Error(), "not found") {
				return &types.Response{
					Code: http.StatusOK,
				}, nil
			} else {
				l.Logger.Error(err)
				return
			}
		}
		json.Unmarshal(rpcResponse.Data, &r)
	} else if ks != nil && ks.IsValid() {
		if r, err = ks.GetPodHorizonAutoscaler(req.Namespace, req.WorkloadName); err != nil {
			if strings.Contains(err.Error(), "could not find") || strings.Contains(err.Error(), "not found") {
				return &types.Response{
					Code: http.StatusOK,
				}, nil
			} else {
				l.Logger.Error(err)
				return
			}
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot GetPodHorizonAutoscaler of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: r,
	}
	return
}
