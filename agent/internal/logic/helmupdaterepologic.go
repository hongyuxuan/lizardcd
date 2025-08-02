package logic

import (
	"context"
	"encoding/json"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type HelmUpdateRepoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	helmService *commonsvc.HelmService
}

func NewHelmUpdateRepoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelmUpdateRepoLogic {
	helmService := commonsvc.NewHelmService(ctx)
	helmService.SetKubeconfig(svcCtx.Config.Kubeconfig)
	return &HelmUpdateRepoLogic{
		ctx:         ctx,
		svcCtx:      svcCtx,
		Logger:      logx.WithContext(ctx),
		helmService: helmService,
	}
}

func (l *HelmUpdateRepoLogic) HelmUpdateRepo(in *agent.HelmEntriesRequest) (*agent.Response, error) {
	var entries []*repo.Entry
	if err := json.Unmarshal(in.Entries, &entries); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := l.helmService.Update(entries); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &agent.Response{
		Code: uint32(codes.OK),
	}, nil
}
