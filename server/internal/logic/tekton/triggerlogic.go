package tekton

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"strings"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/samber/lo"
	"github.com/xdean/goex/xconfig"
	"go.opentelemetry.io/otel"

	"github.com/zeromicro/go-zero/core/logx"
)

type TriggerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTriggerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerLogic {
	return &TriggerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TriggerLogic) Trigger(req *types.TektontriggerReq, secret string) (resp *types.Response, err error) {
	var ciTriggers []commontypes.CiTrigger
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListCiTrigger")).
		Model(&commontypes.CiTrigger{}).
		Where("ci_trigger.git_http_url = ?", req.Project.GitHttpUrl).
		Joins("Application").
		Find(&ciTriggers).Error; err != nil {
		return
	}
	if len(ciTriggers) == 0 {
		return nil, errorx.NewDefaultError("cannot find a webhook trigger")
	}

	var apps []int
	triggerMap := make(map[int]commontypes.CiTrigger)
	for _, hit := range ciTriggers {
		var b []byte
		if b, err = xconfig.Decrypt(hit.Secret[4:], l.svcCtx.Config.Auth.EncKey); err != nil {
			l.Logger.Error(err)
			return
		}
		if string(b) != secret {
			return nil, errorx.NewDefaultError("invalid secret token from HTTP header X-Gitlab-Token")
		}
		triggerMap[hit.AppId] = hit
		for _, path := range hit.TriggerPath {
			for _, commits := range req.Commits {
				// added
				for _, addedPath := range commits.Added {
					if strings.HasPrefix(addedPath, path) {
						apps = append(apps, hit.AppId)
					}
				}
				// modified
				for _, modifiedPath := range commits.Modified {
					if strings.HasPrefix(modifiedPath, path) {
						apps = append(apps, hit.AppId)
					}
				}
				// removed
				for _, removedPath := range commits.Removed {
					if strings.HasPrefix(removedPath, path) {
						apps = append(apps, hit.AppId)
					}
				}
			}
		}
	}
	apps = lo.Uniq(apps)
	for _, app := range apps {
		client := utils.NewHttpClient(otel.Tracer("imroc/req"))
		var tmpl *template.Template
		if tmpl, err = template.New("triggerBody").Parse(triggerMap[app].TriggerBody); err != nil {
			l.Logger.Error(err)
			return
		}
		var buf bytes.Buffer
		if err = tmpl.Execute(&buf, map[string]interface{}{
			"Appname":    triggerMap[app].Application.AppName,
			"Ref":        strings.TrimPrefix(req.Ref, "refs/heads/"),
			"GitSSHUrl":  req.Project.GitSSHUrl,
			"GitHttpUrl": req.Project.GitHttpUrl,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		bodyString := html.UnescapeString(buf.String())
		var body interface{}
		if err = json.Unmarshal([]byte(bodyString), &body); err != nil {
			l.Logger.Error(err)
			return
		}
		if l.svcCtx.Config.Log.Level == "debug" {
			client.EnableDebug(true)
		}
		if err = client.SetBaseURL(triggerMap[app].TriggerEndpoint).
			Post("/").
			SetBody(body).
			Do(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "http.Tektontrigger")).Err; err != nil {
			l.Logger.Error(err)
			return nil, fmt.Errorf("error in send request to http.Tektontrigger: %w", err)
		}
	}
	resp = &types.Response{
		Code: http.StatusOK,
	}
	return
}
