package logic

import (
	"context"
	"fmt"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyTektonResourceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *commonsvc.K8sService
}

func NewApplyTektonResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyTektonResourceLogic {
	return &ApplyTektonResourceLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: commonsvc.NewK8sService(ctx, svcCtx.Clientset, svcCtx.Dynamicclient, svcCtx.TektonClient, svcCtx.TriggerClient),
	}
}

func (l *ApplyTektonResourceLogic) ApplyTektonResource(in *agent.TektonYamlRequest) (*agent.Response, error) {
	l.Logger.Debug("\n", in.Ymlstring)
	if err := l.K8sService.UpdateFromYaml(in.Namespace, in.Ymlstring, ""); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("update YAML failed, error: %v", err))
	}
	return &agent.Response{
		Code:    uint32(codes.OK),
		Message: "update YAML success",
	}, nil
}
