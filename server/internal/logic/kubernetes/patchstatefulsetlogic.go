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

	"github.com/zeromicro/go-zero/core/logx"
)

type PatchStatefulsetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPatchStatefulsetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchStatefulsetLogic {
	return &PatchStatefulsetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PatchStatefulsetLogic) PatchStatefulset(req *types.PatchWorkloadReq) (resp *types.Response, err error) {
	l.Logger.Infof("Patch statefulsets cluster=%s namespace=%s workload=%s container=%s image=%s", req.Cluster, req.Namespace, req.WorkloadName, req.Container, req.Image)
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	if rpcResponse, err = ag.PatchStatefulset(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.PatchStatefulset"), &agent.PatchWorkloadRequest{
		Namespace:    req.Namespace,
		WorkloadName: req.WorkloadName,
		Container:    req.Container,
		Image:        req.Image,
	}); err != nil {
		l.Logger.Error(err)
		return
	}
	var r map[string]interface{}
	json.Unmarshal(rpcResponse.Data, &r)
	resp = &types.Response{
		Code: http.StatusOK,
		Data: r,
	}
	return
}
