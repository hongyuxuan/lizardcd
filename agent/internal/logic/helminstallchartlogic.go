package logic

import (
	"context"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type HelmInstallChartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	helmService *commonsvc.HelmService
}

func NewHelmInstallChartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelmInstallChartLogic {
	helmService := commonsvc.NewHelmService(ctx)
	helmService.SetKubeconfig(svcCtx.Config.Kubeconfig)
	return &HelmInstallChartLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		helmService: helmService,
	}
}

func (l *HelmInstallChartLogic) HelmInstallChart(in *agent.HelmInstallChartRequest) (*agent.Response, error) {
	if err := l.helmService.InstallChart(in.Namespace, in.RepoUrl, in.ChartName, in.ChartVersion, in.ReleaseName, in.Values, in.Wait, time.Duration(in.Timeout)*time.Second); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &agent.Response{
		Code: uint32(codes.OK),
	}, nil
}
