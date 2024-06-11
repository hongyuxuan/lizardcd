package kubernetes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	v1 "k8s.io/api/apps/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeploymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeploymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeploymentLogic {
	return &GetDeploymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeploymentLogic) GetDeployment(req *types.RolloutReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	if rpcResponse, err = ag.GetDeployment(context.WithValue(l.ctx, "SpanName", "rpc.GetDeployment"), &agent.GetWorkloadRequest{
		Namespace:    req.Namespace,
		WorkloadName: req.WorkloadName,
	}); err != nil {
		l.Logger.Error(err)
		return
	}
	var r *v1.Deployment
	json.Unmarshal(rpcResponse.Data, &r)
	resp = &types.Response{
		Code: http.StatusOK,
		Data: r,
	}
	return
}
