package helm

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UninstallChartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUninstallChartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UninstallChartLogic {
	return &UninstallChartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UninstallChartLogic) UninstallChart(req *types.ListReleasesReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	if ag != nil {
		if _, err = ag.HelmUninstallChart(l.ctx, &lizardagent.HelmInstallChartRequest{
			Namespace:   req.Namespace,
			ReleaseName: req.ReleaseName,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
	} else if ks != nil && ks.IsValid() {
		if err = ks.HelmService.UninstallChart(req.Namespace, req.ReleaseName); err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot HelmUninstallChart of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "卸载任务已提交",
	}
	return
}
