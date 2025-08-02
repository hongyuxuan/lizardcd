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

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkloadPodStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkloadPodStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkloadPodStatusLogic {
	return &WorkloadPodStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkloadPodStatusLogic) WorkloadPodStatus(req *types.RolloutReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var r *commontypes.WorkloadStatus
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.GetResourceStatus(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetResourceStatus"), &agent.GetResourceRequest{
			Namespace:    req.Namespace,
			ResourceName: req.WorkloadName,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		json.Unmarshal(rpcResponse.Data, &r)
	} else if ks != nil && ks.IsValid() {
		if r, err = ks.GetResourceStatus(req.Namespace, req.WorkloadType, req.WorkloadName); err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot GetResourceStatus of cluster=%s namespace=%s workload=%s", req.Cluster, req.Namespace, req.WorkloadName)
	}

	resp = &types.Response{
		Code: http.StatusOK,
		Data: r,
	}
	return
}
