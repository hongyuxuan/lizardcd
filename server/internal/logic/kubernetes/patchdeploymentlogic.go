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

type PatchDeploymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPatchDeploymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchDeploymentLogic {
	return &PatchDeploymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PatchDeploymentLogic) PatchDeployment(req *types.PatchWorkloadReq) (resp *types.Response, err error) {
	l.Logger.Infof("Patch deployments cluster=%s namespace=%s workload=%s container=%s image=%s", req.Cluster, req.Namespace, req.WorkloadName, req.Container, req.Image)
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	if rpcResponse, err = ag.PatchDeployment(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.PatchDeployment"), &agent.PatchWorkloadRequest{
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
