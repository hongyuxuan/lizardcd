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

type GetWorkloadQuotaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *commonsvc.K8sService
}

func NewGetWorkloadQuotaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkloadQuotaLogic {
	return &GetWorkloadQuotaLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: commonsvc.NewK8sService(ctx, svcCtx.Clientset, svcCtx.Dynamicclient, svcCtx.TektonClient, svcCtx.TriggerClient),
	}
}

func (l *GetWorkloadQuotaLogic) GetWorkloadQuota(in *agent.GetResourceRequest) (*agent.Response, error) {
	res, err := l.K8sService.GetWorkloadQuota(in.Namespace, in.ResourceType, in.ResourceName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	data, _ := json.Marshal(res)
	return &agent.Response{
		Code: uint32(codes.OK),
		Data: data,
	}, nil
}
