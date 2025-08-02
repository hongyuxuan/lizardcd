package tekton

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/samber/lo"
	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"github.com/xanzy/go-gitlab"
	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/yaml.v2"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	uyaml "k8s.io/apimachinery/pkg/util/yaml"
)

type TriggerpipelinerunLogic struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	gitService      *svc.GitService
	applyTekonLogic *ApplyTektonLogic
	tektonService   *svc.TektonService
}

func NewTriggerpipelinerunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TriggerpipelinerunLogic {
	return &TriggerpipelinerunLogic{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		gitService:      svc.NewGitService(ctx, svcCtx),
		applyTekonLogic: NewApplyTektonLogic(context.Background(), svcCtx),
		tektonService:   svc.NewTektonService(ctx, svcCtx),
	}
}

type Param struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type DefaultTekton struct {
	Cluster string `json:"cluster"`
}

func (l *TriggerpipelinerunLogic) Triggerpipelinerun(req *types.TektontriggerReq) (resp *types.Response, err error) {
	resp = &types.Response{
		Code: http.StatusOK,
	}
	var mergeIID int64
	var revision, committer, commitId string
	var action string
	var ignored bool
	if req.ObjectKind == "merge_request" {
		mergeIID = req.ObjectAttributes.IId
		revision = req.ObjectAttributes.TargetBranch
		committer = req.User.Username
		if req.ObjectAttributes.Action != nil {
			action = *req.ObjectAttributes.Action
			if req.ObjectAttributes.MergeCommitSha != nil {
				commitId = (*req.ObjectAttributes.MergeCommitSha)[0:8]
			} else if req.ObjectAttributes.LastCommit != nil {
				commitId = req.ObjectAttributes.LastCommit.Id[0:8]
			}
			switch action {
			case "merge", "open", "reopen":
				l.Logger.Infof("Received webhook: object_kind=%s git_http_url=%s, revision=%s, committer=%s, commit_id=%s, action=%s, merge_iid=%d", req.ObjectKind, req.Project.GitHttpUrl, revision, committer, commitId, action, mergeIID)
			default:
				ignored = true
			}
		} else {
			ignored = true
		}
	} else if req.ObjectKind == "push" {
		revision = strings.TrimPrefix(req.Ref, "refs/heads/")
		if len(req.Commits) > 0 {
			commitId = req.Commits[len(req.Commits)-1].Id[:8]
		}
		committer = req.Username
		l.Logger.Infof("Received webhook: object_kind=%s git_http_url=%s, revision=%s, committer=%s, commit_id=%s", req.ObjectKind, req.Project.GitHttpUrl, revision, committer, commitId)
	} else {
		ignored = true
	}
	hook := commontypes.LogGitWebhook{
		Timestamp:  time.Now(),
		ObjectKind: req.ObjectKind,
		GitHttpUrl: req.Project.GitHttpUrl,
		MergeIID:   mergeIID,
		Revision:   revision,
		Committer:  committer,
		Action:     action,
		CommitId:   commitId,
		Ignored:    ignored,
		Body:       req.ToJsonString(),
	}
	l.svcCtx.Database.Save(&hook)
	if ignored {
		return
	}
	go func() {
		_, git, project, err := l.gitService.GetGitConnection("admin", req.Project.GitHttpUrl)
		if err != nil {
			l.Logger.Error(err)
			l.svcCtx.Database.Model(&hook).Update("reason", err.Error())
			return
		}
		// get default tekton cluster
		var settings commontypes.Settings
		if err = l.svcCtx.Database.WithContext(context.WithValue(context.Background(), commontypes.TraceIDKey{}, "sqlite.GetDefaultTekton")).
			Model(&commontypes.Settings{}).
			Where("setting_key = ?", "default_tekton").
			First(&settings).Error; err != nil {
			l.Logger.Error(err)
			return
		}
		var defaultTekton DefaultTekton
		if err = json.Unmarshal([]byte(settings.SettingValue), &defaultTekton); err != nil {
			l.Logger.Error(err)
			l.svcCtx.Database.Model(&hook).Update("reason", err.Error())
			return
		}

		// get ci_trigger
		var ciTriggers []commontypes.CiTrigger
		if err = l.svcCtx.Database.WithContext(context.WithValue(context.Background(), commontypes.TraceIDKey{}, "sqlite.ListCiTrigger")).
			Model(&commontypes.CiTrigger{}).
			Where("ci_trigger.git_http_url = ?", req.Project.GitHttpUrl).
			Where("ci_trigger.trigger_type = ?", "pipelinerun").
			Joins("Application").
			Find(&ciTriggers).Error; err != nil {
			l.Logger.Error(err)
			l.svcCtx.Database.Model(&hook).Update("reason", err.Error())
			return
		}
		if len(ciTriggers) == 0 {
			err = errorx.NewDefaultError("Cannot find a webhook trigger of git_http_url=%s, trigger_type=pipelinerun", req.Project.GitHttpUrl)
			l.Logger.Error(err)
			l.svcCtx.Database.Model(&hook).Update("reason", err.Error())
			return
		}

		var fileChanges []string
		annotations := map[string]string{
			"git/objectkind":               req.ObjectKind,
			"git/commit-id":                commitId,
			"phecda.pipeline/triggered-by": committer,
		}
		// if merge_request event, GET /projects/:id/merge_requests/:merge_request_iid/changes
		if req.ObjectKind == "merge_request" {
			if fileChanges, err = l.gitService.GetMergeRequestDiffs(git, project, req.ObjectAttributes.IId); err != nil {
				l.Logger.Error(err)
				l.svcCtx.Database.Model(&hook).Update("reason", err.Error())
				return
			}
			// revision = req.ObjectAttributes.TargetBranch
			annotations["git/merge-url"] = req.ObjectAttributes.Url
		} else { // else push event
			for _, commits := range req.Commits {
				fileChanges = append(fileChanges, commits.Added...)
				fileChanges = append(fileChanges, commits.Modified...)
				fileChanges = append(fileChanges, commits.Removed...)
			}
			if len(req.Commits) > 0 {
				annotations["git/commit-url"] = req.Commits[len(req.Commits)-1].Url
			}
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
		apps = lo.Uniq(apps) // get apps from file changes
		var torunRelease, torunPrecheck []string
		for _, app := range apps {
			reRelease := regexp.MustCompile(triggerMap[app].RefPattern)
			triggerEventRelease := triggerMap[app].TriggerEvent
			if reRelease.Match([]byte(revision)) && strings.Contains(triggerEventRelease, req.ObjectKind) {
				// release pipeline only run push or merged
				if req.ObjectKind == "push" || (req.ObjectKind == "merge_request" && action == "merge") {
					torunRelease = append(torunRelease, app)
				}
			}
			if triggerMap[app].PreCheck.Enable {
				rePrecheck := regexp.MustCompile(triggerMap[app].PreCheck.RefPattern)
				triggerEventPrecheck := triggerMap[app].PreCheck.TriggerEvent
				if rePrecheck.Match([]byte(revision)) && strings.Contains(triggerEventPrecheck, req.ObjectKind) {
					// pre-check pipeline only run push or open/reopen
					if req.ObjectKind == "push" || (req.ObjectKind == "merge_request" && (action == "open" || action == "reopen")) {
						torunPrecheck = append(torunPrecheck, app)
					}
				}
			}
		}
		l.Logger.Infof("Will trigger application: %s for release Pipeline", strings.Join(torunRelease, ","))
		if err = l.applyPipelineRun(torunRelease, defaultTekton, revision, annotations, triggerMap, &hook, git, project, "release"); err != nil {
			l.Logger.Error(err)
			l.svcCtx.Database.Model(&hook).Update("reason", err.Error())
			return
		}
		l.Logger.Infof("Will trigger application: %s for pre-check Pipeline", strings.Join(torunPrecheck, ","))
		if err = l.applyPipelineRun(torunPrecheck, defaultTekton, revision, annotations, triggerMap, &hook, git, project, "pre-check"); err != nil {
			l.Logger.Error(err)
			l.svcCtx.Database.Model(&hook).Update("reason", err.Error())
			return
		}
		torun := append(
			lo.Map(torunRelease, func(item string, _ int) string {
				return item + "(release)"
			}),
			lo.Map(torunPrecheck, func(item string, _ int) string {
				return item + "(pre-check)"
			})...)
		l.svcCtx.Database.Model(&hook).Update("trigger_app", strings.Join(torun, ","))
	}()

	return
}

func (l *TriggerpipelinerunLogic) applyPipelineRun(apps []string, defaultTekton DefaultTekton, revision string, annotations map[string]string, triggerMap map[string]commontypes.CiTrigger, hook *commontypes.LogGitWebhook, git *gitlab.Client, project *gitlab.Project, pipelineType string) (err error) {
	// get tekton_template_pipelinerun, for generate pipelinerun
	var temp commontypes.YamlTemplate
	if err = l.svcCtx.Database.WithContext(context.WithValue(context.Background(), commontypes.TraceIDKey{}, "sqlite.GetTemplates")).
		Model(&commontypes.YamlTemplate{}).
		Where("name = ?", "tekton_template_pipelinerun").
		First(&temp).Error; err != nil {
		return fmt.Errorf("error in GetTemplates: Cannot find template tekton_template_pipelinerun, please create it first")
	}
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(defaultTekton.Cluster, "default"); err != nil {
		return
	}
	for _, app := range apps {
		// get pipeline from git repo
		if scanPath := triggerMap[app].ScanPath; scanPath != "" {
			var filecontents []string
			if filecontents, err = l.gitService.GetFilesContent(git, project, &scanPath, &revision, []string{".tekton.yml", ".tekton.yaml"}, []string{}); err != nil {
				l.svcCtx.Database.Model(hook).Update("reason", err)
				continue
			}
			for _, content := range filecontents {
				l.Logger.Debugf("Fetch yaml from git: %s", content)
				for _, yml := range strings.Split(content, "---") {
					d := uyaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(yml), 4096)
					var unstructureObj *unstructured.Unstructured
					if unstructureObj, err = utils.GetUnstructured(d); err != nil {
						l.svcCtx.Database.Model(hook).Update("reason", err)
						continue
					}
					labels := unstructureObj.GetLabels()
					if labelApp, ok := labels["app"]; ok && labelApp == app {
						if unstructureObj.GetKind() == "PipelineRun" {
							unstructureObj.SetAnnotations(lo.Assign(annotations, unstructureObj.GetAnnotations()))
							// cancel pipelinerun
							if triggerMap[app].UniqInstance {
								l.tektonService.CancelPipelineRun(ag, ks, unstructureObj.GetNamespace(), strings.Join(triggerMap[app].MatchLabels, ","))
							}
						}
						// apply pipeline or pipelinerun to tekton
						newYaml, _ := yaml.Marshal(unstructureObj.Object)
						var res *types.Response
						if res, err = l.applyTekonLogic.ApplyTekton(&types.PatchVariableReq{
							Cluster:   defaultTekton.Cluster,
							Namespace: unstructureObj.GetNamespace(),
							Kind:      unstructureObj.GetKind(),
							Content:   string(newYaml),
							Variables: map[string]interface{}{
								"Revision": revision,
							},
						}); err != nil {
							return fmt.Errorf("error in ApplyTekton: %w", err)
						}
						if res.Code != http.StatusOK {
							return fmt.Errorf("error in ApplyTekton: %s", res.Message)
						}
					}
				}
			}
		} else {
			// get Pipeline from tekton
			labels := strings.Join(triggerMap[app].MatchLabels, ",") + ",type=" + pipelineType
			if labels == "" {
				message := fmt.Sprintf("Cannot trigger Pipeline of %s without labels", app)
				l.Logger.Error(message)
				l.svcCtx.Database.Model(hook).Update("reason", message)
				continue
			}
			r := struct {
				Results []tektonv1.Pipeline `json:"results"`
			}{Results: []tektonv1.Pipeline{}}
			var data []byte
			if ag != nil {
				var rpcResponse *agent.Response
				if rpcResponse, err = ag.ListTektonResource(context.Background(), &lizardagent.TektonListRequest{
					Namespace:     "",
					ResourceType:  "pipelines",
					LabelSelector: labels,
				}); err != nil {
					return fmt.Errorf("error in ListTektonResource: %w", err)
				}
				data = rpcResponse.Data
			} else if ks != nil && ks.IsValid() {
				if data, err = ks.TektonService.ListResource("", "pipelines", labels, "", "", 500); err != nil {
					l.Logger.Error(err)
					return
				}
			} else {
				return errorx.NewDefaultError("Cannot ListTektonResource of tekton")
			}
			json.Unmarshal(data, &r)
			if len(r.Results) == 0 {
				message := fmt.Sprintf("Cannot find Pipeline with labels: %s", labels)
				l.Logger.Error(message)
				l.svcCtx.Database.Model(hook).Update("reason", message)
				continue
			}

			// start pipelinerun, maybe multi pipelines
			for _, pipeline := range r.Results {
				// cancel pipelinerun
				if triggerMap[app].UniqInstance {
					l.tektonService.CancelPipelineRun(ag, ks, pipeline.Namespace, labels)
				}
				if pipelineType == "pre-check" {
					if hook.ObjectKind == "merge_request" {
						annotations["phecda.pipeline/pre-check"] = "mr"
					} else if hook.ObjectKind == "push" {
						annotations["phecda.pipeline/pre-check"] = "push"
					}
				}
				// generate pipelineRun yaml
				variables := map[string]interface{}{
					"Namespace":   pipeline.Namespace,
					"Pipeline":    pipeline.Name,
					"Revision":    revision,
					"Annotations": annotations,
				}
				var res *types.Response
				if res, err = l.applyTekonLogic.ApplyTekton(&types.PatchVariableReq{
					Cluster:   defaultTekton.Cluster,
					Namespace: pipeline.Namespace,
					Kind:      "PipelineRun",
					Content:   temp.Content,
					Variables: variables,
				}); err != nil {
					return fmt.Errorf("error in ApplyTekton: %w", err)
				}
				if res.Code != http.StatusOK {
					return fmt.Errorf("error in ApplyTekton: %s", res.Message)
				}
				l.Logger.Infof("Start pipeline=%s, type=%s success", pipeline.Name, pipelineType)
			}
		}
	}
	return
}
