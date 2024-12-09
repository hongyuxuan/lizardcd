package tekton

import (
	"bytes"
	"context"
	"html"
	"html/template"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyTektonLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyTektonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyTektonLogic {
	return &ApplyTektonLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyTektonLogic) ApplyTekton(req *types.PatchVariableReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	yamlString := req.Content
	if req.Variables != nil {
		var tmpl *template.Template
		if tmpl, err = template.New("yamlTemplates").Parse(req.Content); err != nil {
			l.Logger.Error(err)
			return
		}
		var buf bytes.Buffer
		if err = tmpl.Execute(&buf, req.Variables); err != nil {
			l.Logger.Error(err)
			return
		}
		yamlString = html.UnescapeString(buf.String())
	}
	if _, err = ag.ApplyTektonResource(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.ApplyTekton"), &agent.TektonYamlRequest{
		Namespace:    req.Namespace,
		ResourceType: req.Kind,
		Ymlstring:    yamlString,
	}); err != nil {
		l.Logger.Error(err)
		return &types.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}
	return &types.Response{
		Code: http.StatusOK,
	}, nil
}
