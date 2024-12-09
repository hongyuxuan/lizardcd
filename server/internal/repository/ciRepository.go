package repository

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/samber/lo"
	"github.com/xanzy/go-gitlab"
	"github.com/xdean/goex/xconfig"
	"github.com/zeromicro/go-zero/core/logx"
)

type CiRepository struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	ciTrigger  *commontypes.CiTrigger
	gitRepo    *commontypes.GitRepository
	gitService *svc.GitService
}

func NewCiRepository(ctx context.Context, svcCtx *svc.ServiceContext) *CiRepository {
	return &CiRepository{
		Logger:     logx.WithContext(ctx),
		ctx:        ctx,
		svcCtx:     svcCtx,
		ciTrigger:  &commontypes.CiTrigger{},
		gitRepo:    &commontypes.GitRepository{},
		gitService: svc.NewGitService(ctx, svcCtx),
	}
}

func (r *CiRepository) SaveCiTrigger(body map[string]interface{}) (id interface{}, err error) {
	_, _, tenant, _ := utils.GetPayload(r.ctx)
	b, _ := json.Marshal(body)
	if err = json.Unmarshal(b, &r.ciTrigger); err != nil {
		r.Logger.Error(err)
		return
	}
	if r.ciTrigger.Secret != "" {
		r.ciTrigger.Secret = xconfig.EncryptString(r.ciTrigger.Secret, r.svcCtx.Config.Auth.EncKey)
	}

	var repo *commontypes.GitRepository
	var git *gitlab.Client
	var project *gitlab.Project
	if repo, git, project, err = r.gitService.GetGitConnection(tenant[0], r.ciTrigger.GitHttpUrl); err != nil {
		r.Logger.Error(err)
		return
	}
	// create webhook on gitlab
	if r.ciTrigger.Secret == "" {
		r.ciTrigger.Secret = repo.Secret
	}
	var hooks []*gitlab.ProjectHook
	if hooks, err = r.gitService.GetGitProjectWebhooks(git, project); err != nil {
		r.Logger.Error(err)
		return nil, errors.Unwrap(err)
	}
	webhook := repo.WebhookBaseURL + "/lizardcd/tekton/trigger"
	_, found := lo.Find(hooks, func(item *gitlab.ProjectHook) bool {
		return item.URL == webhook
	})
	if !found {
		var b []byte
		if b, err = xconfig.Decrypt(r.ciTrigger.Secret[4:], r.svcCtx.Config.Auth.EncKey); err != nil {
			r.Logger.Error(err)
			return
		}
		if err = r.gitService.AddGitProjectWebhook(git, project, webhook, string(b)); err != nil {
			r.Logger.Error(err)
			return nil, errors.Unwrap(err)
		}
		r.Logger.Infof("add webhook %s to git project %s success", webhook, r.ciTrigger.GitHttpUrl)
	}
	if err = r.svcCtx.Sqlite.WithContext(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "sqlite.CreateCiTrigger")).Save(&r.ciTrigger).Error; err != nil {
		return
	}
	return r.ciTrigger.Id, nil
}

func (r *CiRepository) ListCiTrigger(req *commontypes.GetDataReq) (resp *types.Response, err error) {
	_, role, tenant, _ := utils.GetPayload(r.ctx)
	var data []commontypes.CiTrigger
	var count int64
	joinTable := "Application"
	tx := r.svcCtx.Sqlite.WithContext(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "sqlite.ListCiTrigger")).Model(commontypes.CiTrigger{})
	utils.SetTx(tx, req, &count, role, tenant, &joinTable)
	if err = tx.Find(&data).Error; err != nil {
		r.Logger.Error(err)
		return
	}
	for i, item := range data {
		if item.Secret == "" {
			continue
		}
		var b []byte
		if b, err = xconfig.Decrypt(item.Secret[4:], r.svcCtx.Config.Auth.EncKey); err != nil {
			r.Logger.Error(err)
			return
		}
		item.Secret = string(b)
		data[i] = item
	}
	return &types.Response{
		Code: http.StatusOK,
		Data: commontypes.ListResult{
			Total:   int(count),
			Results: data,
		},
	}, nil
}

func (r *CiRepository) SaveGitRepository(body map[string]interface{}) (id interface{}, err error) {
	b, _ := json.Marshal(body)
	if err = json.Unmarshal(b, &r.gitRepo); err != nil {
		r.Logger.Error(err)
		return
	}
	r.gitRepo.Secret = xconfig.EncryptString(r.gitRepo.Secret, r.svcCtx.Config.Auth.EncKey)
	return r.gitRepo.Id, r.svcCtx.Sqlite.WithContext(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "sqlite.CreateGitRepository")).Save(&r.gitRepo).Error
}

func (r *CiRepository) ListGitRepository(req *commontypes.GetDataReq) (resp *types.Response, err error) {
	_, role, tenant, _ := utils.GetPayload(r.ctx)
	var data []commontypes.GitRepository
	var count int64
	tx := r.svcCtx.Sqlite.WithContext(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "sqlite.ListGitRepository")).Model(&commontypes.GitRepository{})
	utils.SetTx(tx, req, &count, role, tenant, nil)
	if err = tx.Find(&data).Error; err != nil {
		r.Logger.Error(err)
		return
	}
	for i, item := range data {
		if item.Secret == "" {
			continue
		}
		var b []byte
		if b, err = xconfig.Decrypt(item.Secret[4:], r.svcCtx.Config.Auth.EncKey); err != nil {
			r.Logger.Error(err)
			return
		}
		item.Secret = string(b)
		data[i] = item
	}
	return &types.Response{
		Code: http.StatusOK,
		Data: commontypes.ListResult{
			Total:   int(count),
			Results: data,
		},
	}, nil
}
