package tekton

import (
	"bytes"
	"context"
	"html"
	"html/template"
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

type ApplyTektonLogic struct {
	logx.Logger
	ctx           context.Context
	svcCtx        *svc.ServiceContext
	tektonService *svc.TektonService
}

func NewApplyTektonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyTektonLogic {
	return &ApplyTektonLogic{
		Logger:        logx.WithContext(ctx),
		ctx:           ctx,
		svcCtx:        svcCtx,
		tektonService: svc.NewTektonService(ctx, svcCtx),
	}
}

func (l *ApplyTektonLogic) ApplyTekton(req *types.PatchVariableReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
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

	// if req.Kind == "PipelineRun" { // cancel other pipelineruns
	// 	// parse yaml to unstructured
	// 	unstructureList, err := commonutils.ParseYaml(req.Namespace, yamlString)
	// 	if err != nil {
	// 		l.Logger.Error(err)
	// 		return nil, err
	// 	}
	// 	for _, unstructured := range unstructureList {
	// 		var labels []string
	// 		for k, v := range unstructured.GetLabels() {
	// 			labels = append(labels, k+"="+v)
	// 		}
	// 		l.tektonService.CancelPipelineRun(ag, ks, req.Namespace, strings.Join(labels, ","))
	// 	}
	// }

	if ag != nil {
		if _, err = ag.ApplyYaml(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.ApplyTekton"), &agent.YamlRequest{
			Namespace: req.Namespace,
			Ymlstring: yamlString,
			Kind:      req.Kind,
		}); err != nil {
			l.Logger.Error(err)
			return &types.Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}, nil
		}
	} else if ks != nil && ks.IsValid() {
		if err = ks.UpdateFromYaml(req.Namespace, yamlString, req.Kind); err != nil {
			l.Logger.Error(err)
			return &types.Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}, nil
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot ApplyYaml of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}
	return &types.Response{
		Code: http.StatusOK,
	}, nil
}
