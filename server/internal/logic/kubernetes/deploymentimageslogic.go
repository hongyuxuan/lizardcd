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

type DeploymentImagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeploymentImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeploymentImagesLogic {
	return &DeploymentImagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeploymentImagesLogic) DeploymentImages(req *types.ReplicaReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	if rpcResponse, err = ag.GetDeploymentImages(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetDeploymentReplicas"), &agent.ReplicaRequest{
		Namespace:     req.Namespace,
		Workloads:     req.Workloads,
		InitContainer: req.InitContainer,
	}); err != nil {
		l.Logger.Error(err)
		return
	}
	var r interface{}
	json.Unmarshal(rpcResponse.Data, &r)
	resp = &types.Response{
		Code: http.StatusOK,
		Data: r,
	}
	return
}
