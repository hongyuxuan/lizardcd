package kubernetes

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetYamlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetYamlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetYamlLogic {
	return &GetYamlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetYamlLogic) GetYaml(req *types.ResourceReq) (resp string, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.YamlResponse
	if rpcResponse, err = ag.Getyaml(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetAppsYaml"), &agent.GetYamlRequest{
		Namespace:    req.Namespace,
		ResourceType: req.ResourceType,
		ResourceName: req.ResourceName,
	}); err != nil {
		l.Logger.Error(err)
		return
	}
	resp = rpcResponse.Data
	return
}
