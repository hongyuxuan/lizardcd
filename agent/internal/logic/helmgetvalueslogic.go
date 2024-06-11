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

type HelmGetValuesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	helmUtil *utils.HelmUtil
}

func NewHelmGetValuesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelmGetValuesLogic {
	return &HelmGetValuesLogic{
		ctx:      ctx,
		svcCtx:   svcCtx,
		Logger:   logx.WithContext(ctx),
		helmUtil: utils.NewHelmUtil(ctx),
	}
}

func (l *HelmGetValuesLogic) HelmGetValues(in *agent.ListReleasesRequest) (*agent.Response, error) {
	output, err := l.helmUtil.GetValues(in.Namespace, l.svcCtx.Config.Kubeconfig, in.ReleaseName, int(in.Revision))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: output,
	}, nil
}
