package kubernetes

import (
	"bytes"
	"context"
	"encoding/json"
	"html"
	"net/http"
	"text/template"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PatchVariableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPatchVariableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchVariableLogic {
	return &PatchVariableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PatchVariableLogic) PatchVariable(req *types.PatchVariableReq) (resp *types.Response, err error) {
	var tmpl *template.Template
	if tmpl, err = template.New("appTemplates").Parse(req.Content); err != nil {
		l.Logger.Error(err)
		return
	}
	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, req.Variables); err != nil {
		l.Logger.Error(err)
		return
	}
	yamlString := html.UnescapeString(buf.String())

	l.Logger.Infof("apply yaml in cluster=%s namespace=%s variables=%v :\n%s", req.Cluster, req.Namespace, req.Variables, yamlString)

	var ks *commonsvc.K8sService
	var ag lizardagent.LizardAgent
	var data []byte
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.ApplyYaml(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.PatchYaml"), &agent.YamlRequest{
			Namespace: req.Namespace,
			Ymlstring: yamlString,
			Kind:      req.Kind,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		data = rpcResponse.Data
	} else if ks != nil && ks.IsValid() {
		if err = ks.UpdateFromYaml(req.Namespace, yamlString, req.Kind); err != nil {
			l.Logger.Error(err)
			return nil, errorx.NewDefaultError("update YAML failed, error: %v", err)
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot ApplyYaml of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}

	var r map[string]interface{}
	json.Unmarshal(data, &r)
	resp = &types.Response{
		Code: http.StatusOK,
		Data: r,
	}
	return
}
