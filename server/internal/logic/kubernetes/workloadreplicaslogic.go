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

type WorkloadReplicasLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkloadReplicasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkloadReplicasLogic {
	return &WorkloadReplicasLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkloadReplicasLogic) WorkloadReplicas(req *types.ReplicaReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var data []commontypes.WorkloadReplica
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.GetWorkloadReplicas(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetWorkloadReplicas"), &agent.ReplicaRequest{
			Namespace:    req.Namespace,
			WorkloadType: req.WorkloadType,
			Workloads:    req.Workloads,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		json.Unmarshal(rpcResponse.Data, &data)
	} else if ks != nil && ks.IsValid() {
		if data, err = ks.GetReplicas(req.Namespace, req.WorkloadType, req.Workloads); err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot GetWorkloadReplicas of cluster=%s namespace=%s workload=%s", req.Cluster, req.Namespace, req.Workloads)
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: map[string]interface{}{
			"workloads": data,
		},
	}
	return
}
