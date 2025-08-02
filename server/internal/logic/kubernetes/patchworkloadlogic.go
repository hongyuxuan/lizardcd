package kubernetes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PatchWorkloadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPatchWorkloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchWorkloadLogic {
	return &PatchWorkloadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PatchWorkloadLogic) PatchWorkload(req *types.PatchWorkloadReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var data []byte
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.PatchWorkload(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.PatchWorkload"), &agent.PatchWorkloadRequest{
			Namespace:    req.Namespace,
			WorkloadType: req.WorkloadType,
			WorkloadName: req.WorkloadName,
			Container:    req.Container,
			Image:        req.Image,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		data = rpcResponse.Data
	} else if ks != nil && ks.IsValid() {
		if data, err = ks.PatchWorkload(req.Namespace, req.WorkloadType, req.WorkloadName, req.Container, req.Image); err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot PatchWorkload of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}
	var r map[string]interface{}
	json.Unmarshal(data, &r)
	resp = &types.Response{
		Code: http.StatusOK,
		Data: r,
	}
	return
}
