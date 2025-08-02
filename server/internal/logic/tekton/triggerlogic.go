package tekton

import (
	"bytes"
	"context"
	"encoding/json"
	"html"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"

	"github.com/zeromicro/go-zero/core/logx"
)

type TriggerLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	gitService *svc.GitService
}

func NewTriggerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerLogic {
	return &TriggerLogic{
		Logger:     logx.WithContext(ctx),
		ctx:        ctx,
		svcCtx:     svcCtx,
		gitService: svc.NewGitService(ctx, svcCtx),
	}
}

func (l *TriggerLogic) Trigger(req *types.TektontriggerReq) (resp *types.Response, err error) {
	resp = &types.Response{
		Code: http.StatusOK,
	}
	if req.ObjectKind == "merge_request" {
		l.Logger.Infof("Received webhook: object_kind=%s git_http_url=%s, target=%s, author=%s", req.ObjectKind, req.Project.GitHttpUrl, req.ObjectAttributes.TargetBranch, req.User.Username)
	} else if req.ObjectKind == "push" {
		l.Logger.Infof("Received webhook: object_kind=%s git_http_url=%s, revision=%s, author=%s", req.ObjectKind, req.Project.GitHttpUrl, req.Ref, req.Username)
	} else {
		l.Logger.Infof("Received webhook: object_kind=%s ignored", req.ObjectKind)
	}

	go func() {
		var ciTriggers []commontypes.CiTrigger
		if err = l.svcCtx.Database.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListCiTrigger")).
			Model(&commontypes.CiTrigger{}).
			Where("ci_trigger.git_http_url = ?", req.Project.GitHttpUrl).
			Where("ci_trigger.trigger_type != ?", "pipelinerun").
			Joins("Application").
			Find(&ciTriggers).Error; err != nil {
			l.Logger.Error(err)
			return
		}
		if len(ciTriggers) == 0 {
			l.Logger.Errorf("Cannot find a webhook trigger of git_http_url=%s", req.Project.GitHttpUrl)
			return
		}

		var fileChanges []string
		var revision string
		// if merge_request event, GET /projects/:id/merge_requests/:merge_request_iid/changes
		if req.ObjectKind == "merge_request" {
			_, git, project, err := l.gitService.GetGitConnection("admin", req.Project.GitHttpUrl)
			if err != nil {
				l.Logger.Error(err)
				return
			}
			if fileChanges, err = l.gitService.GetMergeRequestDiffs(git, project, req.ObjectAttributes.IId); err != nil {
				l.Logger.Error(err)
				return
			}
			revision = req.ObjectAttributes.TargetBranch
		} else { // else push event
			for _, commits := range req.Commits {
				fileChanges = append(fileChanges, commits.Added...)
				fileChanges = append(fileChanges, commits.Modified...)
				fileChanges = append(fileChanges, commits.Removed...)
			}
			revision = strings.TrimPrefix(req.Ref, "refs/heads/")
		}
		var apps []string
		triggerMap := make(map[string]commontypes.CiTrigger)
		for _, hit := range ciTriggers {
			triggerMap[hit.Application.AppName] = hit
			for _, path := range hit.TriggerPath {
				for _, file := range fileChanges {
					if strings.HasPrefix(file, path) {
						apps = append(apps, hit.Application.AppName)
					}
				}
			}
		}
		apps = lo.Uniq(apps)
		torun := lo.Filter(apps, func(app string, _ int) bool {
			re := regexp.MustCompile(triggerMap[app].RefPattern)
			triggerEvent := triggerMap[app].TriggerEvent
			return re.Match([]byte(revision)) && strings.Contains(triggerEvent, req.ObjectKind)
		})
		for _, app := range torun {
			client := utils.NewHttpClient(otel.Tracer("imroc/req"))
			var tmpl *template.Template
			var triggerBody string
			if triggerMap[app].TriggerBody != nil {
				triggerBody = *triggerMap[app].TriggerBody
			}
			if tmpl, err = template.New("triggerBody").Parse(triggerBody); err != nil {
				l.Logger.Error(err)
				return
			}
			var buf bytes.Buffer
			if err = tmpl.Execute(&buf, map[string]interface{}{
				"Appname":    triggerMap[app].Application.AppName,
				"Ref":        revision,
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
				return
			}
		}
		l.Logger.Infof("Triggered application: %s", strings.Join(torun, ","))
	}()
	return
}
