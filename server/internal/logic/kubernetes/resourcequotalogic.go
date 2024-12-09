package kubernetes

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/constant"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"

	corev1 "k8s.io/api/core/v1"
)

type ResourceQuotaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResourceQuotaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceQuotaLogic {
	return &ResourceQuotaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceQuotaLogic) ResourceQuota(req *types.ResourceReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	var rr *corev1.ResourceRequirements
	var resourceList []map[string]string
	if req.ResourceType == constant.K8S_RESOURCE_TYPE_DEPLOYMENTS {
		for _, d := range strings.Split(req.ResourceName, ",") {
			if rpcResponse, err = ag.GetDeploymentQuota(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetDeploymentQuota"), &agent.GetWorkloadRequest{
				Namespace:    req.Namespace,
				WorkloadName: d,
			}); err != nil {
				l.Logger.Error(err)
				continue
			}
			json.Unmarshal(rpcResponse.Data, &rr)
			resourceRes := map[string]string{
				"name":            d,
				"limits_cpu":      rr.Limits.Cpu().String(),
				"limits_memory":   rr.Limits.Memory().String(),
				"requests_cpu":    rr.Requests.Cpu().String(),
				"requests_memory": rr.Requests.Memory().String(),
			}
			resourceList = append(resourceList, resourceRes)
		}
	}
	if req.ResourceType == constant.K8S_RESOURCE_TYPE_STATEFULSETS {
		for _, d := range strings.Split(req.ResourceName, ",") {
			if rpcResponse, err = ag.DeleteStatefulset(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetStatefulsetQuota"), &agent.GetWorkloadRequest{
				Namespace:    req.Namespace,
				WorkloadName: d,
			}); err != nil {
				l.Logger.Error(err)
				continue
			}
			json.Unmarshal(rpcResponse.Data, &rr)
			resourceRes := map[string]string{
				"name":            d,
				"limits.cpu":      rr.Limits.Cpu().String(),
				"limits.memory":   rr.Limits.Memory().String(),
				"requests.cpu":    rr.Requests.Cpu().String(),
				"requests.memory": rr.Requests.Memory().String(),
			}
			resourceList = append(resourceList, resourceRes)
		}
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: resourceList,
	}
	return
}
