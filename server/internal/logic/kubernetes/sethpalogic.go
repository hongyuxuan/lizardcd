package kubernetes

import (
	"context"
	"encoding/json"
	"net/http"

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

type SethpaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSethpaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SethpaLogic {
	return &SethpaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SethpaLogic) Sethpa(req *types.HpaReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var data *autov2.HorizontalPodAutoscaler
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.SetPodHorizonAutoscaler(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.SetPodHorizonAutoscaler"), &agent.HpaRequest{
			Namespace:    req.Namespace,
			WorkloadName: req.WorkloadName,
			Max:          req.Max,
			Min:          req.Min,
			Cpu:          req.Cpu,
			Memory:       req.Memory,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		json.Unmarshal(rpcResponse.Data, &data)
	} else if ks != nil && ks.IsValid() {
		if data, err = ks.SetPodHorizonAutoscaler(req.Namespace, req.WorkloadName, &req.Max, &req.Min, &req.Cpu, &req.Memory); err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot SetPodHorizonAutoscaler of cluster=%s namespace=%s workload=%s", req.Cluster, req.Namespace, req.WorkloadName)
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "设置HPA成功",
		Data:    data,
	}
	return
}
