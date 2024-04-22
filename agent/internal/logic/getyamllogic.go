package logic

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"google.golang.org/grpc/codes"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetyamlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *svc.K8sService
}

func NewGetyamlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetyamlLogic {
	return &GetyamlLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: svc.GetK8sService(ctx, svcCtx),
	}
}

func (l *GetyamlLogic) Getyaml(in *agent.GetYamlRequest) (resp *agent.YamlResponse, err error) {
	var res string
	if in.ResourceType == "deployments" || in.ResourceType == "statefulsets" {
		if res, err = l.K8sService.GetAppsV1ResourceYAML(in.Namespace, in.ResourceType, in.ResourceName); err != nil {
			l.Logger.Error(err)
			return
		}
	} else if in.ResourceType == "ingresses" {
		if res, err = l.K8sService.GetIngressYAML(in.Namespace, in.ResourceName); err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		if res, err = l.K8sService.GetCoreV1ResourceYAML(in.Namespace, in.ResourceType, in.ResourceName); err != nil {
			l.Logger.Error(err)
			return
		}
	}
	resp = &agent.YamlResponse{
		Code: uint32(codes.OK),
		Data: res,
	}
	return
}
