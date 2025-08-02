package kubernetes

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScaleWorkloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScaleWorkloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScaleWorkloadLogic {
	return &ScaleWorkloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScaleWorkloadLogic) ScaleWorkload(req *types.ScaleReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	for _, workload := range req.Workloads {
		if workload.Disabled {
			continue
		}
		if ag != nil {
			if _, err = ag.ScaleWorkload(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.ScaleWorkload"), &agent.ScaleRequest{
				Namespace:    req.Namespace,
				WorkloadType: req.WorkloadType,
				WorkloadName: workload.Name,
				Replicas:     uint32(workload.Replicas),
			}); err != nil {
				l.Logger.Error(err)
				return
			}
		} else if ks != nil && ks.IsValid() {
			if err = ks.ScaleWorkload(req.Namespace, req.WorkloadType, workload.Name, uint32(workload.Replicas)); err != nil {
				l.Logger.Error(err)
				return
			}
		} else {
			return nil, errorx.NewDefaultError("Cannot RolloutWorkload of cluster=%s namespace=%s", req.Cluster, req.Namespace)
		}
	}
	resp = &types.Response{
		Code: http.StatusOK,
	}
	return
}
