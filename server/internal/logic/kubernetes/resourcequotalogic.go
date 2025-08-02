package kubernetes

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"

	v1 "k8s.io/api/core/v1"
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
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var resourceList []map[string]string
	var data *v1.ResourceRequirements
	for _, d := range strings.Split(req.ResourceName, ",") {
		if ag != nil {
			var rpcResponse *agent.Response
			if rpcResponse, err = ag.GetWorkloadQuota(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetWorkloadQuota"), &agent.GetResourceRequest{
				Namespace:    req.Namespace,
				ResourceType: req.ResourceType,
				ResourceName: d,
			}); err != nil {
				l.Logger.Error(err)
				continue
			}
			json.Unmarshal(rpcResponse.Data, &data)
		} else if ks != nil && ks.IsValid() {
			if data, err = ks.GetWorkloadQuota(req.Namespace, req.ResourceType, d); err != nil {
				l.Logger.Error(err)
				return
			}
		} else {
			return nil, errorx.NewDefaultError("Cannot GetWorkloadQuota of cluster=%s namespace=%s workload=%s", req.Cluster, req.Namespace, d)
		}
		resourceRes := map[string]string{
			"name":            d,
			"limits_cpu":      data.Limits.Cpu().String(),
			"limits_memory":   data.Limits.Memory().String(),
			"requests_cpu":    data.Requests.Cpu().String(),
			"requests_memory": data.Requests.Memory().String(),
		}
		resourceList = append(resourceList, resourceRes)
	}

	resp = &types.Response{
		Code: http.StatusOK,
		Data: resourceList,
	}
	return
}
