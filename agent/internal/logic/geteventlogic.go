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

type GetEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *svc.K8sService
}

func NewGetEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventLogic {
	return &GetEventLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: svc.GetK8sService(ctx, svcCtx),
	}
}

func (l *GetEventLogic) GetEvent(in *agent.GetEventRequest) (*agent.Response, error) {
	res, err := l.K8sService.GetEvents(in.Namespace, in.ObjectKind, in.ObjectName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data, _ := json.Marshal(res)
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
