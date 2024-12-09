package svc

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hongyuxuan/lizardcd/common/constant"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/xanzy/go-gitlab"
	"github.com/zeromicro/go-zero/core/logx"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	uyaml "k8s.io/apimachinery/pkg/util/yaml"
)

type GitService struct {
	logx.Logger
	ctx    context.Context
	svcCtx *ServiceContext
}

func NewGitService(ctx context.Context, svcCtx *ServiceContext) *GitService {
	return &GitService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (g *GitService) GetGitConnection(tenant, gitHttpUrl string) (repo *commontypes.GitRepository, git *gitlab.Client, project *gitlab.Project, err error) {
	var gitRepos []commontypes.GitRepository
	tx := g.svcCtx.Sqlite.WithContext(context.WithValue(context.Background(), commontypes.TraceIDKey{}, "sqlite.ListGitRepository")).Model(&commontypes.GitRepository{})
	if tenant != constant.ROLE_ADMIN {
		tx.Where("tenant = ?", tenant)
	}
	if err = tx.Find(&gitRepos).Error; err != nil {
		g.Logger.Error(err)
		return
	}
	for _, gitrepo := range gitRepos {
		if strings.HasPrefix(gitHttpUrl, gitrepo.GitHttpUrl) {
			if git, err = gitlab.NewClient(gitrepo.AccessToken, gitlab.WithBaseURL(fmt.Sprintf("%s/api/v4", gitrepo.GitHttpUrl))); err != nil {
				return nil, nil, nil, fmt.Errorf("error in gitlab.NewClient: %w", err)
			}
			var project *gitlab.Project
			if project, err = getGitProjectByHttpUrl(git, gitHttpUrl); err != nil {
				return nil, nil, nil, fmt.Errorf("error in getGitProjectByHttpUrl: %w", err)
			}
			return &gitrepo, git, project, nil
		}
	}
	return nil, nil, nil, errorx.NewDefaultError("cannot find git settings for %s", gitHttpUrl)
}

func (g *GitService) GetGitProjectWebhooks(git *gitlab.Client, project *gitlab.Project) (hooks []*gitlab.ProjectHook, err error) {
	if hooks, _, err = git.Projects.ListProjectHooks(project.ID, &gitlab.ListProjectHooksOptions{}); err != nil {
		return nil, fmt.Errorf("error in git.Projects.ListProjectHooks: %w", err)
	}
	return
}

func (g *GitService) AddGitProjectWebhook(git *gitlab.Client, project *gitlab.Project, webhookUrl, secret string) (err error) {
	push := true
	if _, _, err = git.Projects.AddProjectHook(project.ID, &gitlab.AddProjectHookOptions{
		URL:        &webhookUrl,
		PushEvents: &push,
		Token:      &secret,
	}); err != nil {
		return fmt.Errorf("error in git.Projects.AddProjectHook: %w", err)
	}
	return
}

func (g *GitService) GetYamlFromGit(application commontypes.Application, tenant string) (commitId, filecontent string, err error) {
	var git *gitlab.Client
	var project *gitlab.Project
	if _, git, project, err = g.GetGitConnection(tenant, application.GitHttpUrl); err != nil {
		return "", "", fmt.Errorf("error in GetGitConnection: %w", err)
	}
	var commit *gitlab.Commit
	if commit, err = g.GetLastCommitId(git, project, &application.GitOps.Path, &application.GitOps.GitRevision); err != nil {
		return
	}
	commitId = commit.ShortID
	g.Logger.Infof("get project %s success", project.WebURL)
	includes := strings.Split(application.GitOps.Include, ",")
	excludes := strings.Split(application.GitOps.Exclude, ",")
	filecontent, err = g.getFileContent(git, project, &application.GitOps.Path, &application.GitOps.GitRevision, includes, excludes)
	return
}

func (g *GitService) getFileContent(git *gitlab.Client, project *gitlab.Project, path, ref *string, includes, excludes []string) (filecontent string, err error) {
	recursive := false
	page := 1
	perpage := 20
	var candidates []string
	for {
		var nodes []*gitlab.TreeNode
		if nodes, _, err = git.Repositories.ListTree(project.ID, &gitlab.ListTreeOptions{
			Path:        path,
			Ref:         ref,
			Recursive:   &recursive,
			ListOptions: gitlab.ListOptions{Page: page, PerPage: perpage}}); err != nil {
			return "", fmt.Errorf("error in git.Repositories.ListTree: %w", err)
		}
		if len(nodes) == 0 {
			break
		}
		for _, node := range nodes {
			hit := false
			for _, include := range includes {
				re := regexp.MustCompile(include)
				if re.Match([]byte(node.Name)) {
					hit = true
				}
			}
			for _, exclude := range excludes {
				if exclude == "" {
					continue
				}
				re := regexp.MustCompile(exclude)
				if re.Match([]byte(node.Name)) {
					hit = false
				}
			}
			if hit {
				g.Logger.Infof("hit %s ref=%s file=%s", project.WebURL, *ref, node.Path)
				candidates = append(candidates, node.Path)
			}
		}
		for _, candidate := range candidates {
			var b []byte
			if b, _, err = git.RepositoryFiles.GetRawFile(project.ID, candidate, &gitlab.GetRawFileOptions{Ref: ref}); err != nil {
				return "", fmt.Errorf("error in getFileContent: %w", err)
			}
			filecontent += string(b) + "\n---\n"
		}
		g.Logger.Debugf("\n%s", filecontent)
		page += 1
	}
	return
}

func (g *GitService) GetLastCommitId(git *gitlab.Client, project *gitlab.Project, path, ref *string) (commit *gitlab.Commit, err error) {
	var commits []*gitlab.Commit
	if commits, _, err = git.Commits.ListCommits(project.ID, &gitlab.ListCommitsOptions{
		RefName:     ref,
		Path:        path,
		ListOptions: gitlab.ListOptions{PerPage: 1, Page: 1},
	}); err != nil {
		return nil, fmt.Errorf("error in git.Commits.ListCommits: %w", err)
	}
	if len(commits) > 0 {
		return commits[0], nil
	}
	return
}

func (g *GitService) ParseYaml(ctx context.Context, application commontypes.Application, yamlstring string) (err error) {
	// clear application_resource
	g.svcCtx.Sqlite.WithContext(context.WithValue(ctx, commontypes.TraceIDKey{}, "sqlite.DeleteApplicationResource")).Delete(&commontypes.ApplicationResource{}, "application_id = ?", application.Id)

	for _, yml := range strings.Split(yamlstring, "---") {
		if strings.TrimSpace(yml) == "" {
			continue
		}
		d := uyaml.NewYAMLOrJSONDecoder(bytes.NewBufferString(yml), 4096)
		var unstructureObj *unstructured.Unstructured
		if unstructureObj, err = utils.GetUnstructured(d); err != nil {
			return fmt.Errorf("error in GetUnstructured: %w", err)
		}
		for _, workload := range application.Workload {
			appResource := commontypes.ApplicationResource{
				ApplicationId:   application.Id,
				Cluster:         workload.Cluster,
				Namespace:       workload.Namespace,
				ResourceType:    strings.ToLower(unstructureObj.GetKind()) + "s",
				ResourceName:    unstructureObj.GetName(),
				DesiredManifest: yml,
			}
			if err = g.svcCtx.Sqlite.WithContext(context.WithValue(ctx, commontypes.TraceIDKey{}, "sqlite.SaveApplicationResource")).Create(&appResource).Error; err != nil {
				g.Logger.Error(err)
				return fmt.Errorf("error in SaveApplicationResource: %w", err)
			}
			g.Logger.Debugf("Saved application=%s cluster=%s namespace=%s resourceType=%s resourceName=%s to database success", application.AppName, workload.Cluster, workload.Namespace, unstructureObj.GetKind(), unstructureObj.GetName())
		}
	}
	return
}

func getGitProjectByHttpUrl(git *gitlab.Client, gitHttpUrl string) (project *gitlab.Project, err error) {
	page := 1
	perpage := 20
	projectName, found := strings.CutSuffix(filepath.Base(gitHttpUrl), ".git")
	if !found {
		return nil, errorx.NewDefaultError("git url must be end with .git")
	}
	for {
		var projects []*gitlab.Project
		if projects, _, err = git.Projects.ListProjects(&gitlab.ListProjectsOptions{
			Search:      &projectName,
			ListOptions: gitlab.ListOptions{Page: page, PerPage: perpage}}); err != nil {
			return nil, fmt.Errorf("error in git.Projects.ListProjects: %w", err)
		}
		if len(projects) == 0 {
			break
		}
		for _, hit := range projects {
			if hit.HTTPURLToRepo == gitHttpUrl {
				project = hit
				break
			}
		}
		page += 1
	}
	if project == nil {
		err = errorx.NewDefaultError("cannot find project %s", gitHttpUrl)
	}
	return
}
