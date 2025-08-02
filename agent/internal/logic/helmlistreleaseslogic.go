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

type HelmListReleasesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	helmService *commonsvc.HelmService
}

func NewHelmListReleasesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelmListReleasesLogic {
	helmService := commonsvc.NewHelmService(ctx)
	helmService.SetKubeconfig(svcCtx.Config.Kubeconfig)
	return &HelmListReleasesLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		helmService: helmService,
	}
}

func (l *HelmListReleasesLogic) HelmListReleases(in *agent.ListReleasesRequest) (*agent.Response, error) {
	res, err := l.helmService.ListRelease(in.Namespace, in.ReleaseName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data, _ := json.Marshal(res)
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
