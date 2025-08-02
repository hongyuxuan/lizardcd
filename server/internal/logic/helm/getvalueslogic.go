package helm

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
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
	var ks *commonsvc.K8sService
	var r []byte
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.HelmGetValues(l.ctx, &lizardagent.ListReleasesRequest{
			Namespace:   req.Namespace,
			ReleaseName: req.ReleaseName,
			Revision:    req.Revision,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		r = rpcResponse.Data
	} else if ks != nil && ks.IsValid() {
		if r, err = ks.HelmService.GetValues(req.Namespace, req.ReleaseName, int(req.Revision)); err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		return "", errorx.NewDefaultError("Cannot HelmGetValues of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}
	content = string(r)
	return
}
