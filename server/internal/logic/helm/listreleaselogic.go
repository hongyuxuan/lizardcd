package helm

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

	"github.com/zeromicro/go-zero/core/logx"
)

type ListReleaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListReleaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListReleaseLogic {
	return &ListReleaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListReleaseLogic) ListRelease(req *types.ListReleasesReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var r []commontypes.ReleaseElement
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.HelmListReleases(l.ctx, &agent.ListReleasesRequest{
			Namespace:   req.Namespace,
			ReleaseName: req.ReleaseName,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		json.Unmarshal(rpcResponse.Data, &r)
	} else if ks != nil && ks.IsValid() {
		if r, err = ks.HelmService.ListRelease(req.Namespace, req.ReleaseName); err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot HelmListReleases of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}
	if r == nil {
		r = make([]commontypes.ReleaseElement, 0)
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: r,
	}
	return
}
