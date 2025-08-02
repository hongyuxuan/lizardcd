package kubernetes

import (
	"context"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetYamlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetYamlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetYamlLogic {
	return &GetYamlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetYamlLogic) GetYaml(req *types.ResourceReq) (resp string, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	if ag != nil {
		var rpcResponse *agent.YamlResponse
		if rpcResponse, err = ag.GetYaml(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetResourceYaml"), &agent.GetYamlRequest{
			Namespace:    req.Namespace,
			ResourceType: req.ResourceType,
			ResourceName: req.ResourceName,
		}); err != nil {
			l.Logger.Error(err)
			return err.Error(), nil
		}
		return rpcResponse.Data, nil
	} else if ks != nil && ks.IsValid() {
		if resp, err = ks.GetResourceYAML(req.Namespace, req.ResourceType, req.ResourceName); err != nil {
			l.Logger.Error(err)
			return err.Error(), nil
		}
		return
	} else {
		return errorx.NewDefaultError("Cannot GetYaml of cluster=%s namespace=%s", req.Cluster, req.Namespace).Error(), nil
	}
}
