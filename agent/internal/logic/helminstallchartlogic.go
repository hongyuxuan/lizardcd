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

type HelmInstallChartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	helmUtil *utils.HelmUtil
}

func NewHelmInstallChartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelmInstallChartLogic {
	return &HelmInstallChartLogic{
		ctx:      ctx,
		svcCtx:   svcCtx,
		Logger:   logx.WithContext(ctx),
		helmUtil: utils.NewHelmUtil(ctx),
	}
}

func (l *HelmInstallChartLogic) HelmInstallChart(in *agent.HelmInstallChartRequest) (*agent.Response, error) {
	if err := l.helmUtil.InstallChart(in.Namespace, l.svcCtx.Config.Kubeconfig, in.RepoUrl, in.ChartName, in.ChartVersion, in.ReleaseName, in.Values, in.Wait, time.Duration(in.Timeout)*time.Second); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &agent.Response{
		Code: uint32(codes.OK),
	}, nil
}
