package logic

import (
	"context"
	"encoding/json"
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
		if len(strings.TrimSpace(y)) == 0 {
			continue
		}
		taskResult := make(chan map[string]interface{})
		chArr = append(chArr, taskResult)
		go l.K8sService.UpdateFromYaml(in.Namespace, y, in.Kind, taskResult)
	}
	failed := []string{}
	var workloadType, workloadName string
	firstWorkload := true
	for _, ch := range chArr {
		res := <-ch
		if res["success"] == false {
			failed = append(failed, res["message"].(string))
		} else {
			if (res["workloadType"].(string) == "Deployment" ||
				res["workloadType"].(string) == "StatefulSet" ||
				res["workloadType"].(string) == "DaemonSet") && firstWorkload { // 返回yaml里的第一个工作负载
				workloadType = res["workloadType"].(string)
				workloadName = res["workloadName"].(string)
				firstWorkload = false
			}
		}
	}
	if len(failed) > 0 {
		return nil, status.Error(codes.Internal, fmt.Sprintf("update YAML failed, errmsg: %v", failed))
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
