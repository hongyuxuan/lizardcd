package logic

import (
	"context"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type HelmUpgradeChartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	helmUtil *utils.HelmUtil
}

func NewHelmUpgradeChartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelmUpgradeChartLogic {
	return &HelmUpgradeChartLogic{
		ctx:      ctx,
		svcCtx:   svcCtx,
		Logger:   logx.WithContext(ctx),
		helmUtil: utils.NewHelmUtil(ctx),
	}
}

func (l *HelmUpgradeChartLogic) HelmUpgradeChart(in *agent.HelmInstallChartRequest) (*agent.Response, error) {
	err := l.helmUtil.UpgradeChart(in.Namespace, l.svcCtx.Config.Kubeconfig, in.RepoUrl, in.ReleaseName, in.ChartName, in.ChartVersion, int(in.Revision), in.Values, in.Wait, time.Duration(in.Timeout)*time.Second)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &agent.Response{
		Code: uint32(codes.OK),
	}, nil
}
