package db

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListdataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListdataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListdataLogic {
	return &ListdataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListdataLogic) Listdata(req *commontypes.GetDataReq) (resp *types.Response, err error) {
	if req.Tablename == "application" {
		var data []commontypes.Application
		return l.list(l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListApplication")).Model(commontypes.Application{}), data, req)
	} else if req.Tablename == "user" {
		var data []commontypes.User
		tx := l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListUser")).Model(&commontypes.User{}).Select("id", "username", "role", "tenant", "update_at")
		return l.list(tx, data, req)
	} else if req.Tablename == "task_history" {
		var data []commontypes.TaskHistory
		tx := l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListTaskHistory")).Model(commontypes.TaskHistory{})
		if req.Preload {
			tx.Preload("TaskHistoryWorkloads")
		}
		return l.list(tx, data, req)
	} else if req.Tablename == "helm_repositories" {
		var data []commontypes.HelmRepositories
		return l.list(l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListHelmRepositories")).Model(commontypes.HelmRepositories{}), data, req)
	} else {
		// []map[string]interface{} cannot use l.list, will be failed with 'sql: Scan error on column index 0, name "id": destination not a pointer'
		_, role, tenant, _ := utils.GetPayload(l.ctx)
		var data []map[string]interface{}
		tx := l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListData")).Table(req.Tablename)
		var count int64
		utils.SetTx(tx, &count, req, role, tenant)
		if err = tx.Find(&data).Error; err != nil {
			l.Logger.Error(err)
			return
		}
		if data == nil {
			data = []map[string]interface{}{}
		}
		resp = &types.Response{
			Code: http.StatusOK,
			Data: map[string]interface{}{
				"total":   count,
				"results": data,
			},
		}
		return
	}
}

func (l *ListdataLogic) list(tx *gorm.DB, models any, req *commontypes.GetDataReq) (resp *types.Response, err error) {
	_, role, tenant, _ := utils.GetPayload(l.ctx)
	var count int64
	utils.SetTx(tx, &count, req, role, tenant)
	if err = tx.Find(&models).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: map[string]interface{}{
			"total":   count,
			"results": models,
		},
	}
	return
}
