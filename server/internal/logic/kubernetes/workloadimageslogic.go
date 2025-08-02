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

type WorkloadImagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWorkloadImagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WorkloadImagesLogic {
	return &WorkloadImagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WorkloadImagesLogic) WorkloadImages(req *types.ReplicaReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var data []commontypes.WorkloadImage
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.GetWorkloadImages(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetWorkloadImages"), &agent.ReplicaRequest{
			Namespace:     req.Namespace,
			WorkloadType:  req.WorkloadType,
			Workloads:     req.Workloads,
			InitContainer: req.InitContainer,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		json.Unmarshal(rpcResponse.Data, &data)
	} else if ks != nil && ks.IsValid() {
		if data, err = ks.GetImages(req.Namespace, req.WorkloadType, req.InitContainer, req.Workloads); err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot GetWorkloadImages of cluster=%s namespace=%s workloads=%v", req.Cluster, req.Namespace, req.Workloads)
	}

	resp = &types.Response{
		Code: http.StatusOK,
		Data: data,
	}
	return
}
