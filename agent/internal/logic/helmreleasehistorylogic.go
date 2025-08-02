package logic

import (
	"context"
	"encoding/json"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type HelmReleaseHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	helmService *commonsvc.HelmService
}

func NewHelmReleaseHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelmReleaseHistoryLogic {
	helmService := commonsvc.NewHelmService(ctx)
	helmService.SetKubeconfig(svcCtx.Config.Kubeconfig)
	return &HelmReleaseHistoryLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		helmService: helmService,
	}
}

func (l *HelmReleaseHistoryLogic) HelmReleaseHistory(in *agent.ListReleasesRequest) (*agent.Response, error) {
	res, err := l.helmService.ReleaseHistory(in.Namespace, in.ReleaseName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data, _ := json.Marshal(res)
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
