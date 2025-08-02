package logic

import (
	"context"
	"encoding/json"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type DockerCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDockerCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DockerCheckLogic {
	return &DockerCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DockerCheckLogic) DockerCheck(in *agent.DockerDeployRequest) (*agent.Response, error) {
	containerInfo, err := l.svcCtx.DockerClient.ContainerInspect(l.ctx, in.ContainerName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	b, _ := json.Marshal(containerInfo.State)
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: b,
	}, nil
}
