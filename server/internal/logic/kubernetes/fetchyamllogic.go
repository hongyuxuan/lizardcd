package kubernetes

import (
	"bytes"
	"context"
	"html"
	"net/http"
	"text/template"

	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FetchYamlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFetchYamlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FetchYamlLogic {
	return &FetchYamlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FetchYamlLogic) FetchYaml(req *types.FetchYamlReq) (resp *types.Response, err error) {
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
	return &types.Response{
		Code: http.StatusOK,
		Data: yamlString,
	}, nil
}
