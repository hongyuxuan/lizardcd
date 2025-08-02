package logic

import (
	"context"
	"encoding/json"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyYamlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *commonsvc.K8sService
}

func NewApplyYamlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyYamlLogic {
	return &ApplyYamlLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: commonsvc.NewK8sService(ctx, svcCtx.Clientset, svcCtx.Dynamicclient, svcCtx.TektonClient, svcCtx.TriggerClient),
	}
}

func (l *ApplyYamlLogic) ApplyYaml(in *agent.YamlRequest) (resp *agent.Response, err error) {
	var workloadType, workloadName string
	firstWorkload := true
	unstructureList, err := utils.ParseYaml(in.Namespace, in.Ymlstring)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Parse YAML failed, error: %v", err)
	}
	for _, unstructureObj := range unstructureList {
		resourceName := unstructureObj.GetName()
		if resourceName == "" {
			resourceName = unstructureObj.GetGenerateName()
		}
		resourceType := unstructureObj.GetKind()
		if (resourceType == "Deployment" || resourceType == "StatefulSet" || resourceType == "DaemonSet") && firstWorkload {
			workloadType = resourceType
			workloadName = resourceName
			firstWorkload = false
		}
	}
	if err := l.K8sService.UpdateFromYaml(in.Namespace, in.Ymlstring, in.Kind); err != nil {
		return nil, status.Errorf(codes.Internal, "Update YAML failed, error: %v", err)
	}
	data := struct {
		WorkloadType string
		WorkloadName string
	}{
		WorkloadType: workloadType,
		WorkloadName: workloadName,
	}
	b, _ := json.Marshal(data)
	return &agent.Response{
		Code:    uint32(codes.OK),
		Message: "update YAML success",
		Data:    b,
	}, nil
}
