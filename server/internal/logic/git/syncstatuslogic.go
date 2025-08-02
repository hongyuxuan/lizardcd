package git

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/xanzy/go-gitlab"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncstatusLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	gitService  *svc.GitService
	taskService *svc.TaskService
}

func NewSyncstatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncstatusLogic {
	return &SyncstatusLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		gitService:  svc.NewGitService(ctx, svcCtx),
		taskService: svc.NewTaskService(ctx, svcCtx),
	}
}

func (l *SyncstatusLogic) Syncstatus(req *types.SyncStatusReq) (resp *types.Response, err error) {
	_, _, tenant, _ := utils.GetPayload(l.ctx)
	data := make(map[string]*types.SyncStatus)
	apps := strings.Split(req.Apps, ",")
	var wg sync.WaitGroup
	resultChan := make(chan *types.SyncStatus, len(apps))
	for _, app := range apps {
		wg.Add(1)
		go l.getSyncStatus(app, tenant[0], resultChan, &wg)
	}
	wg.Wait()
	close(resultChan)
	for result := range resultChan {
		if result.Error != nil {
			return nil, result.Error
		}
		if !result.Failed {
			data[result.AppName] = result
		}
	}
	return &types.Response{
		Code: http.StatusOK,
		Data: data,
	}, nil
}

func (l *SyncstatusLogic) getSyncStatus(appName, tenant string, ch chan *types.SyncStatus, wg *sync.WaitGroup) {
	var application commontypes.Application
	if err := l.svcCtx.Database.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.GetApplication")).
		First(&application, "app_name = ?", appName).Error; err != nil {
		l.Logger.Error(err)
		ch <- &types.SyncStatus{Failed: true}
		wg.Done()
		return
	}
	var resData *types.SyncStatus
	// 获取gitlab同步状态
	_, git, project, err := l.gitService.GetGitConnection(tenant, application.GitHttpUrl)
	if err != nil {
		l.Logger.Debug(err)
		ch <- &types.SyncStatus{Failed: true}
		wg.Done()
		return
	}
	if application.GitOps == nil {
		l.Logger.Debugf("GitOps of application=%s not setting", appName)
		ch <- &types.SyncStatus{Failed: true}
		wg.Done()
		return
	}
	var commit *gitlab.Commit
	if commit, err = l.gitService.GetLastCommitId(git, project, &application.GitOps.Path, &application.GitOps.GitRevision); err != nil {
		l.Logger.Error(err)
		ch <- &types.SyncStatus{Failed: true, Error: err}
		wg.Done()
		return
	}
	resData = &types.SyncStatus{
		AppName:        appName,
		LocalCommitId:  application.GitOps.CommitId,
		RemoteCommitId: commit.ShortID,
		Comment:        commit.Message,
		Author:         commit.CommitterName,
		Status:         "已同步",
	}
	if commit.ShortID != application.GitOps.CommitId {
		resData.Status = "未同步"
		resData.Reason = fmt.Sprintf("当前应用 commitId=%s 与最新的 commitId=%s 不一致", application.GitOps.CommitId, commit.ShortID)
		ch <- resData
		wg.Done()
		return
	}
	// 获取工作负载同步状态
	var appResource []commontypes.ApplicationResource
	if err = l.svcCtx.Database.Model(&commontypes.ApplicationResource{}).
		Where("application_id = ?", application.Id).
		Find(&appResource).Error; err != nil {
		l.Logger.Error(err)
		ch <- &types.SyncStatus{Failed: true, Error: err}
		wg.Done()
		return
	}
	for _, resource := range appResource {
		manifests, err := l.taskService.GetResourceManifests(l.ctx, resource)
		if err != nil {
			l.Logger.Error(err)
			resData.Status = "Error"
			resData.Reason = errors.Unwrap(err).Error()
			ch <- resData
			wg.Done()
			return
		}
		if resource.LastManifest != manifests {
			resData.Status = "未同步"
			resData.Reason = fmt.Sprintf("%s=%s manifests is different from current K8S", resource.ResourceType, resource.ResourceName)
			ch <- resData
			wg.Done()
			return
		}
	}
	ch <- resData
	wg.Done()
}
