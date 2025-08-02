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
	tektontypes "github.com/hongyuxuan/tekton-sdk-go/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PatchTektonResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPatchTektonResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchTektonResourceLogic {
	return &PatchTektonResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PatchTektonResourceLogic) PatchTektonResource(req *types.ResourceReq, data []tektontypes.PatchOptions) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var resData []byte
	b, _ := json.Marshal(data)
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.PatchTektonResource(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.PatchTektonResource"), &agent.TektonPatchRequest{
			Namespace:    req.Namespace,
			ResourceType: req.ResourceType,
			ResourceName: req.ResourceName,
			PatchOptions: b,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		resData = rpcResponse.Data
	} else if ks != nil && ks.IsValid() {
		if resData, err = ks.TektonService.PatchResource(req.Namespace, req.ResourceType, req.ResourceName, b); err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot PatchTektonResource of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}
	var res interface{}
	json.Unmarshal(resData, &res)
	return &types.Response{
		Code: http.StatusOK,
		Data: res,
	}, nil
}
