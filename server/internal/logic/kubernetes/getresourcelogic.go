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

type GetResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetResourceLogic {
	return &GetResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetResourceLogic) GetResource(req *types.ResourceReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var data []byte
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.GetResource(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.ListDeployment"), &agent.GetResourceRequest{
			Namespace:    req.Namespace,
			ResourceType: req.ResourceType,
			ResourceName: req.ResourceName,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		data = rpcResponse.Data
	} else if ks != nil && ks.IsValid() {
		if data, _, err = ks.GetResource(req.Namespace, req.ResourceType, req.ResourceName); err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot GetResource of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}
	var r interface{}
	json.Unmarshal(data, &r)
	return &types.Response{
		Code: http.StatusOK,
		Data: r,
	}, nil
}
