package tekton

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

type GetTektonResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTektonResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTektonResourceLogic {
	return &GetTektonResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTektonResourceLogic) GetTektonResource(req *types.ResourceReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
	}
	var data []byte
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.GetTektonResource(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.ListTektonResource"), &agent.TektonYamlRequest{
			Namespace:    req.Namespace,
			ResourceType: req.ResourceType,
			ResourceName: req.ResourceName,
		}); err != nil {
			l.Logger.Error(err)
			return &types.Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}, nil
		}
		data = rpcResponse.Data
	} else if ks != nil && ks.IsValid() {
		if data, err = ks.TektonService.GetResource(req.Namespace, req.ResourceType, req.ResourceName); err != nil {
			l.Logger.Error(err)
			return &types.Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}, nil
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot GetTektonResource of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}
	var res interface{}
	json.Unmarshal(data, &res)
	resp.Data = res
	return
}
