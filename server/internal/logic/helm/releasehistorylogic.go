package helm

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

type ReleaseHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReleaseHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReleaseHistoryLogic {
	return &ReleaseHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReleaseHistoryLogic) ReleaseHistory(req *types.ListReleasesReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	if rpcResponse, err = ag.HelmReleaseHistory(l.ctx, &agent.ListReleasesRequest{
		Namespace:   req.Namespace,
		ReleaseName: req.ReleaseName,
	}); err != nil {
		l.Logger.Error(err)
		return
	}
	var r []commontypes.ReleaseHistoryInfo
	json.Unmarshal(rpcResponse.Data, &r)
	if r == nil {
		r = make([]commontypes.ReleaseHistoryInfo, 0)
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: r,
	}
	return
}
