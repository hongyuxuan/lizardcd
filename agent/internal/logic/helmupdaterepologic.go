package logic

import (
	"context"
	"encoding/json"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type HelmUpdateRepoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	helmUtil *utils.HelmUtil
}

func NewHelmUpdateRepoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelmUpdateRepoLogic {
	return &HelmUpdateRepoLogic{
		ctx:      ctx,
		svcCtx:   svcCtx,
		Logger:   logx.WithContext(ctx),
		helmUtil: utils.NewHelmUtil(ctx),
	}
}

func (l *HelmUpdateRepoLogic) HelmUpdateRepo(in *agent.HelmEntriesRequest) (*agent.Response, error) {
	var entries []*repo.Entry
	if err := json.Unmarshal(in.Entries, &entries); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err := l.helmUtil.Update(entries); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &agent.Response{
		Code: uint32(codes.OK),
	}, nil
}
