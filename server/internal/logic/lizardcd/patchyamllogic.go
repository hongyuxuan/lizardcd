package lizardcd

import (
	"bytes"
	"context"
	"encoding/json"
	"html"
	"net/http"
	"text/template"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
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

func (l *PatchYamlLogic) PatchYaml(req *types.PatchYamlReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
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
	var rpcResponse *agent.Response
	if rpcResponse, err = ag.ApplyYaml(context.WithValue(l.ctx, "SpanName", "rpc.PatchYaml"), &agent.YamlRequest{
		Namespace: req.Namespace,
		Ymlstring: yamlString,
	}); err != nil {
		l.Logger.Error(err)
		return
	}
	var r map[string]interface{}
	json.Unmarshal(rpcResponse.Data, &r)
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: rpcResponse.Message,
	}
	return
}
