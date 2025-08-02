package tekton

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

type GetTektonYamlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTektonYamlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTektonYamlLogic {
	return &GetTektonYamlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTektonYamlLogic) GetTektonYaml(req *types.ResourceReq) (resp string, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	if ag != nil {
		var rpcResponse *agent.YamlResponse
		if rpcResponse, err = ag.GetTektonYaml(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetTektonYaml"), &agent.TektonYamlRequest{
			Namespace:    req.Namespace,
			ResourceType: req.ResourceType,
			ResourceName: req.ResourceName,
			WithStatus:   req.WithStatus,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		return rpcResponse.Data, nil
	} else if ks != nil && ks.IsValid() {
		if resp, err = ks.TektonService.GetResourceYAML(req.Namespace, req.ResourceType, req.ResourceName, req.WithStatus); err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		return "", errorx.NewDefaultError("Cannot GetTektonYaml of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}
	return
}
