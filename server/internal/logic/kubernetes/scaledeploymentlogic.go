package kubernetes

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScaleDeploymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScaleDeploymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScaleDeploymentLogic {
	return &ScaleDeploymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScaleDeploymentLogic) ScaleDeployment(req *types.ScaleReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	for _, workload := range req.Workloads {
		if workload.Disabled {
			continue
		}
		if _, err = ag.ScaleDeployment(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.ScaleDeployment"), &agent.ScaleRequest{
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
