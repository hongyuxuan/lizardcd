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

type PatchYamlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPatchYamlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchYamlLogic {
	return &PatchYamlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PatchYamlLogic) PatchYaml(req *types.PatchYamlReq, body string) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var data []byte
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.ApplyYaml(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.PatchYaml"), &agent.YamlRequest{
			Namespace: req.Namespace,
			Ymlstring: body,
			Kind:      req.Kind,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		data = rpcResponse.Data
	} else if ks != nil && ks.IsValid() {
		if err = ks.UpdateFromYaml(req.Namespace, body, req.Kind); err != nil {
			l.Logger.Error(err)
			return nil, errorx.NewDefaultError("update YAML failed, error: %v", err)
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot ApplyYaml of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}

	var r map[string]interface{}
	json.Unmarshal(data, &r)
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "update YAML success",
		Data:    r,
	}
	return
}
