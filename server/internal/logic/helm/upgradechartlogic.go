package helm

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpgradeChartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpgradeChartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpgradeChartLogic {
	return &UpgradeChartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpgradeChartLogic) UpgradeChart(req *types.InstallChartReq) (resp *types.Response, err error) {
	_, _, tenant, _ := utils.GetPayload(l.ctx)
	wait, timeout, err := l.svcCtx.GetHelmSettings(tenant)
	if err != nil {
		l.Logger.Error(err)
		return
	}
	// create a context with 2s timeout, if there is no error within 2s, return "安装任务已提交"
	ctx, cancel := context.WithTimeout(l.ctx, 2*time.Second)
	defer cancel()
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	if _, err = ag.HelmUpgradeChart(ctx, &lizardagent.HelmInstallChartRequest{
		RepoUrl:      req.RepoUrl,
		Namespace:    req.Namespace,
		ChartName:    req.ChartName,
		ChartVersion: req.ChartVersion,
		ReleaseName:  req.ReleaseName,
		Revision:     req.Revision,
		Values:       []byte(req.Values),
		Wait:         wait,
		Timeout:      timeout,
	}); err != nil {
		if strings.Contains(err.Error(), "DeadlineExceeded") { // timeout because --wait
			resp = &types.Response{
				Code:    http.StatusOK,
				Message: "重装任务已提交",
			}
			return resp, nil
		}
		l.Logger.Error(err)
		return
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "重装任务已提交",
	}
	return
}
