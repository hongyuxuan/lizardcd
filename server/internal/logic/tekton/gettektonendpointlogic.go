package tekton

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTektonEndpointLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTektonEndpointLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTektonEndpointLogic {
	return &GetTektonEndpointLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTektonEndpointLogic) GetTektonEndpoint(req *types.PatchYamlReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, _, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var endpoint string
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.GetTektonEndpoint(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.ListTektonResource"), &agent.TektonEndpointRequest{}); err != nil {
			l.Logger.Error(err)
			return
		}
		endpoint = string(rpcResponse.Data)
	} else {
		endpoint = l.svcCtx.Config.TektonEndpoint
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: endpoint,
	}
	return
}
