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
	corev1 "k8s.io/api/core/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatefulsetPodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatefulsetPodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatefulsetPodsLogic {
	return &StatefulsetPodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatefulsetPodsLogic) StatefulsetPods(req *types.RolloutReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	if rpcResponse, err = ag.GetStatefulsetPod(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.StatefulsetPods"), &agent.GetWorkloadRequest{
		Namespace:    req.Namespace,
		WorkloadName: req.WorkloadName,
	}); err != nil {
		l.Logger.Error(err)
		return
	}
	var r []corev1.Pod
	json.Unmarshal(rpcResponse.Data, &r)
	resp = &types.Response{
		Code: http.StatusOK,
		Data: r,
	}
	return
}
