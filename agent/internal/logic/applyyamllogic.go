package logic

import (
	"context"
	"fmt"
	"strings"

	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyYamlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	K8sService *svc.K8sService
}

func NewApplyYamlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyYamlLogic {
	return &ApplyYamlLogic{
		ctx:        ctx,
		svcCtx:     svcCtx,
		Logger:     logx.WithContext(ctx),
		K8sService: svc.GetK8sService(ctx, svcCtx),
	}
}

func (l *ApplyYamlLogic) ApplyYaml(in *agent.YamlRequest) (*agent.Response, error) {
	yamlArr := strings.Split(in.Ymlstring, "---")
	var chArr []chan map[string]interface{}
	for _, y := range yamlArr {
		taskResult := make(chan map[string]interface{})
		chArr = append(chArr, taskResult)
		go l.K8sService.UpdateFromYaml(in.Namespace, y, in.Kind, taskResult)
	}
	failed := []string{}
	for _, ch := range chArr {
		res := <-ch
		if res["success"] == false {
			failed = append(failed, res["message"].(string))
		}
	}
	if len(failed) > 0 {
		return nil, status.Error(codes.Internal, fmt.Sprintf("update YAML failed, errmsg: %v", failed))
	}
	return &agent.Response{
		Code:    uint32(codes.OK),
		Message: "update YAML success",
	}, nil
}
