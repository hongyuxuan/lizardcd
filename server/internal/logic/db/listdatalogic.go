package db

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/constant"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/repository"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListdataLogic struct {
	logx.Logger
	ctx           context.Context
	svcCtx        *svc.ServiceContext
	ciRepository  *repository.CiRepository
	appRepository *repository.ApplicationRepository
}

func NewListdataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListdataLogic {
	return &ListdataLogic{
		Logger:        logx.WithContext(ctx),
		ctx:           ctx,
		svcCtx:        svcCtx,
		ciRepository:  repository.NewCiRepository(ctx, svcCtx),
		appRepository: repository.NewApplicationRepository(ctx, svcCtx),
	}
}

func (l *ListdataLogic) Listdata(req *commontypes.GetDataReq) (resp *types.Response, err error) {
	_, role, tenant, _ := utils.GetPayload(l.ctx)
	switch req.Tablename {
	case "application":
		var data []commontypes.Application
		joinTable := "Template"
		return l.list(l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListApplication")).Model(commontypes.Application{}), data, req, &joinTable)
	case "application_faas":
		var data []commontypes.ApplicationFaas
		return l.list(l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListApplicationFaas")).Model(commontypes.ApplicationFaas{}), data, req, nil)
	case "application_resource":
		return l.appRepository.ListResource(req)
	case "user":
		var data []commontypes.User
		tx := l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListUser")).Model(&commontypes.User{}).Select("id", "username", "role", "tenant", "update_at")
		return l.list(tx, data, req, nil)
	case "tokens":
		if role != constant.ROLE_ADMIN {
			return nil, errorx.NewError(http.StatusForbidden, "not permitted", nil)
		}
		var data []commontypes.Tokens
		return l.list(l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListTokens")).Model(&commontypes.Tokens{}), data, req, nil)
	case "task_history":
		var data []commontypes.TaskHistory
		tx := l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListTaskHistory")).Model(&commontypes.TaskHistory{})
		if req.Preload {
			tx.Preload("TaskHistoryWorkloads")
		}
		return l.list(tx, data, req, nil)
	case "helm_repositories":
		var data []commontypes.HelmRepositories
		return l.list(l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListHelmRepositories")).Model(commontypes.HelmRepositories{}), data, req, nil)
	case "git_repository":
		return l.ciRepository.ListGitRepository(req)
	case "ci_trigger":
		return l.ciRepository.ListCiTrigger(req)
	default:
		// []map[string]interface{} cannot use l.list, will be failed with 'sql: Scan error on column index 0, name "id": destination not a pointer'
		data := []map[string]interface{}{}
		tx := l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListData")).Table(req.Tablename)
		var count int64
		utils.SetTx(tx, req, &count, role, tenant, nil)
		if err = tx.Find(&data).Error; err != nil {
			l.Logger.Error(err)
			return
		}
		resp = &types.Response{
			Code: http.StatusOK,
			Data: commontypes.ListResult{
				Total:   int(count),
				Results: data,
			},
		}
		return
	}
}

func (l *ListdataLogic) list(tx *gorm.DB, models any, req *commontypes.GetDataReq, joinTable *string) (resp *types.Response, err error) {
	_, role, tenant, _ := utils.GetPayload(l.ctx)
	var count int64
	utils.SetTx(tx, req, &count, role, tenant, joinTable)
	if err = tx.Find(&models).Error; err != nil {
		l.Logger.Error(err)
		return
	}

	resp = &types.Response{
		Code: http.StatusOK,
		Data: commontypes.ListResult{
			Total:   int(count),
			Results: models,
		},
	}
	return
}
