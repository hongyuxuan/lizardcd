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

type PatchConfigmapLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPatchConfigmapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchConfigmapLogic {
	return &PatchConfigmapLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PatchConfigmapLogic) PatchConfigmap(req *types.PatchConfigReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}

	var data []byte
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.PatchConfigmap(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.PatchConfigmap"), &agent.PatchConfigmapRequest{
			Namespace:    req.Namespace,
			ResourceType: req.ResourceType,
			ResourceName: req.ResourceName,
			Key:          req.Key,
			Value:        req.Value,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		data = rpcResponse.Data
	} else if ks != nil && ks.IsValid() {
		if data, err = ks.PatchConfig(req.Namespace, req.ResourceType, req.ResourceName, req.Key, req.Value); err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot PatchConfig of cluster=%s namespace=%s %s=%s", req.Cluster, req.Namespace, req.ResourceType, req.ResourceName)
	}
	var res interface{}
	json.Unmarshal(data, &res)
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "更新成功",
		Data:    res,
	}
	return
}
