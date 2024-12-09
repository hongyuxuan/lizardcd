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

type RolloutDeploymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRolloutDeploymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RolloutDeploymentLogic {
	return &RolloutDeploymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RolloutDeploymentLogic) RolloutDeployment(req *types.RolloutReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	if rpcResponse, err = ag.RolloutDeployment(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.RolloutDeployment"), &agent.RolloutWorkloadRequest{
		Namespace:    req.Namespace,
		WorkloadName: req.WorkloadName,
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
