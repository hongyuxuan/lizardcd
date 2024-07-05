package kubernetes

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/constant"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteResourceLogic {
	return &DeleteResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteResourceLogic) DeleteResource(req *types.ResourceReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	if req.ResourceType == constant.K8S_RESOURCE_TYPE_DEPLOYMENTS {
		if _, err = ag.DeleteDeployment(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.DeleteDeployment"), &agent.GetWorkloadRequest{
			Namespace:    req.Namespace,
			WorkloadName: req.ResourceName,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
	}
	if req.ResourceType == constant.K8S_RESOURCE_TYPE_STATEFULSETS {
		if _, err = ag.DeleteStatefulset(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.DeleteStatefulset"), &agent.GetWorkloadRequest{
			Namespace:    req.Namespace,
			WorkloadName: req.ResourceName,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
	}
	if req.ResourceType == constant.ISTIO_CRD_TYPE_DESTINATIONRULE {
		if _, err = ag.DeleteDestinationRule(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.DeleteDestinationRule"), &agent.IstioGetRequest{
			Namespace: req.Namespace,
			Name:      req.ResourceName,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
	}
	if req.ResourceType == constant.ISTIO_CRD_TYPE_VIRTUALSERVICE {
		if _, err = ag.DeleteVirtualService(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.DeleteVirtualService"), &agent.IstioGetRequest{
			Namespace: req.Namespace,
			Name:      req.ResourceName,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
	}
	resp = &types.Response{
		Code: http.StatusOK,
	}
	return
}
