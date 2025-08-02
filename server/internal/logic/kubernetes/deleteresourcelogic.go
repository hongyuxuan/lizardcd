package kubernetes

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/constant"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
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
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	switch req.ResourceType {
	case constant.ISTIO_CRD_TYPE_DESTINATIONRULE:
		if _, err = ag.DeleteDestinationRule(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.DeleteDestinationRule"), &agent.IstioGetRequest{
			Namespace: req.Namespace,
			Name:      req.ResourceName,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
	case constant.ISTIO_CRD_TYPE_VIRTUALSERVICE:
		if _, err = ag.DeleteVirtualService(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.DeleteVirtualService"), &agent.IstioGetRequest{
			Namespace: req.Namespace,
			Name:      req.ResourceName,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
	case constant.ISTIO_CRD_TYPE_GATEWAY:
		if _, err = ag.DeleteGateway(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.DeleteGateway"), &agent.IstioGetRequest{
			Namespace: req.Namespace,
			Name:      req.ResourceName,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
	default:
		if ag != nil {
			if _, err = ag.DeleteResource(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.DeleteDeployment"), &agent.DeleteResourceRequest{
				Namespace:    req.Namespace,
				ResourceType: req.ResourceType,
				ResourceName: req.ResourceName,
				Force:        req.Force,
			}); err != nil {
				l.Logger.Error(err)
				return
			}
		} else if ks != nil && ks.IsValid() {
			if err = ks.DeleteResource(req.Namespace, req.ResourceType, req.ResourceName, req.Force); err != nil {
				l.Logger.Error(err)
				return
			}
		} else {
			return nil, errorx.NewDefaultError("Cannot DeleteResource of cluster=%s namespace=%s", req.Cluster, req.Namespace)
		}
	}
	resp = &types.Response{
		Code: http.StatusOK,
	}
	return
}
