package helm

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetValuesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetValuesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetValuesLogic {
	return &GetValuesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetValuesLogic) GetValues(req *types.ListReleasesReq) (content string, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	if rpcResponse, err = ag.HelmGetValues(l.ctx, &lizardagent.ListReleasesRequest{
		Namespace:   req.Namespace,
		ReleaseName: req.ReleaseName,
		Revision:    req.Revision,
	}); err != nil {
		l.Logger.Error(err)
		return
	}
	content = string(rpcResponse.Data)
	return
}
