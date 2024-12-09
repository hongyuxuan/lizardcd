package logic

import (
	"context"
	"encoding/json"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type HelmReleaseHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	helmUtil *utils.HelmUtil
}

func NewHelmReleaseHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelmReleaseHistoryLogic {
	return &HelmReleaseHistoryLogic{
		ctx:      ctx,
		svcCtx:   svcCtx,
		Logger:   logx.WithContext(ctx),
		helmUtil: utils.NewHelmUtil(ctx),
	}
}

func (l *HelmReleaseHistoryLogic) HelmReleaseHistory(in *agent.ListReleasesRequest) (*agent.Response, error) {
	res, err := l.helmUtil.ReleaseHistory(in.Namespace, l.svcCtx.Config.Kubeconfig, in.ReleaseName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data, _ := json.Marshal(res)
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
