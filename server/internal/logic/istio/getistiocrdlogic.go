package istio

import (
	"context"
	"encoding/json"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"gopkg.in/yaml.v2"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetIstioCrdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetIstioCrdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIstioCrdLogic {
	return &GetIstioCrdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIstioCrdLogic) GetIstioCrd(req *types.ResourceReq) (resp string, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	var r map[string]interface{}
	if req.ResourceType == "destinationrules" {
		if rpcResponse, err = ag.GetDestinationRule(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetDestinationRule"), &agent.IstioGetRequest{
			Namespace: req.Namespace,
			Name:      req.ResourceName,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
	} else if req.ResourceType == "virtualservices" {
		if rpcResponse, err = ag.GetVirtualService(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetVirtualService"), &agent.IstioGetRequest{
			Namespace: req.Namespace,
			Name:      req.ResourceName,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
	}
	json.Unmarshal(rpcResponse.Data, &r)
	delete(r, "status")
	b, _ := yaml.Marshal(r)
	return string(b), nil
}
