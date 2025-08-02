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

type HelmUpgradeChartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	helmService *commonsvc.HelmService
}

func NewHelmUpgradeChartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelmUpgradeChartLogic {
	helmService := commonsvc.NewHelmService(ctx)
	helmService.SetKubeconfig(svcCtx.Config.Kubeconfig)
	return &HelmUpgradeChartLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		helmService: helmService,
	}
}

func (l *HelmUpgradeChartLogic) HelmUpgradeChart(in *agent.HelmInstallChartRequest) (*agent.Response, error) {
	err := l.helmService.UpgradeChart(in.Namespace, in.RepoUrl, in.ReleaseName, in.ChartName, in.ChartVersion, int(in.Revision), in.Values, in.Wait, time.Duration(in.Timeout)*time.Second)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &agent.Response{
		Code: uint32(codes.OK),
	}, nil
}
