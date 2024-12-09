package logic

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type HelmUninstallChartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	helmUtil *utils.HelmUtil
}

func NewHelmUninstallChartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelmUninstallChartLogic {
	return &HelmUninstallChartLogic{
		ctx:      ctx,
		svcCtx:   svcCtx,
		Logger:   logx.WithContext(ctx),
		helmUtil: utils.NewHelmUtil(ctx),
	}
}

func (l *HelmUninstallChartLogic) HelmUninstallChart(in *agent.HelmInstallChartRequest) (*agent.Response, error) {
	if err := l.helmUtil.UninstallChart(in.Namespace, l.svcCtx.Config.Kubeconfig, in.ReleaseName); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &agent.Response{
		Code: uint32(codes.OK),
	}, nil
}
