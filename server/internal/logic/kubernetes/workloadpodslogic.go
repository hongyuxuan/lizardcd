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
	corev1 "k8s.io/api/core/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type WorkloadPodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkloadPodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkloadPodsLogic {
	return &WorkloadPodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkloadPodsLogic) WorkloadPods(req *types.RolloutReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var r []corev1.Pod
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.GetWorkloadPod(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetWorkloadPod"), &agent.GetResourceRequest{
			Namespace:    req.Namespace,
			ResourceType: req.WorkloadType,
			ResourceName: req.WorkloadName,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		json.Unmarshal(rpcResponse.Data, &r)
	} else if ks != nil && ks.IsValid() {
		_, labels, err := ks.GetResource(req.Namespace, req.WorkloadType, req.WorkloadName)
		if err != nil {
			l.Logger.Error(err)
			return nil, err
		}
		if r, _, err = ks.ListPods(req.Namespace, labels, "", "", 0); err != nil {
			l.Logger.Error(err)
			return nil, err
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot GetWorkloadPods of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}

	resp = &types.Response{
		Code: http.StatusOK,
		Data: r,
	}
	return
}
