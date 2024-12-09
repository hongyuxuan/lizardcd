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

type InstallChartLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInstallChartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InstallChartLogic {
	return &InstallChartLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InstallChartLogic) InstallChart(req *types.InstallChartReq) (resp *types.Response, err error) {
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
	if _, err = ag.HelmInstallChart(ctx, &lizardagent.HelmInstallChartRequest{
		RepoUrl:      req.RepoUrl,
		ChartName:    req.ChartName,
		ChartVersion: req.ChartVersion,
		Namespace:    req.Namespace,
		ReleaseName:  req.ReleaseName,
		Values:       []byte(req.Values),
		Wait:         wait,
		Timeout:      timeout,
	}); err != nil {
		if strings.Contains(err.Error(), "DeadlineExceeded") { // timeout because --wait
			resp = &types.Response{
				Code:    http.StatusOK,
				Message: "安装任务已提交",
			}
			return resp, nil
		}
		l.Logger.Error(err)
		return
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "安装任务已提交",
	}
	return
}
