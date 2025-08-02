package db

import (
	"context"
	"errors"
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/constant"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/repository"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetdataLogic struct {
	logx.Logger
	ctx           context.Context
	svcCtx        *svc.ServiceContext
	appRepository *repository.ApplicationRepository
}

func NewGetdataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetdataLogic {
	return &GetdataLogic{
		Logger:        logx.WithContext(ctx),
		ctx:           ctx,
		svcCtx:        svcCtx,
		appRepository: repository.NewApplicationRepository(ctx, svcCtx),
	}
}

func (l *GetdataLogic) Getdata(req *types.DataByIdReq) (resp *types.Response, err error) {
	_, role, tenant, _ := utils.GetPayload(l.ctx)
	switch req.Tablename {
	case "application":
		var application *commontypes.Application
		if application, err = l.appRepository.Get(req.Id, role, tenant); err != nil {
			l.Logger.Error(err)
			return nil, errors.Unwrap(err)
		}
		resp = &types.Response{
			Code: http.StatusOK,
			Data: application,
		}
		break
	case "task_history":
		var taskHistory commontypes.TaskHistory
		tx := l.svcCtx.Database.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.GetTaskHistory"))
		if role != constant.ROLE_ADMIN {
			tx.Where("tenant IN ?", tenant)
		}
		if err = tx.Preload("TaskHistoryWorkloads").First(&taskHistory, "id = ?", req.Id).Error; err != nil {
			l.Logger.Error(err)
			return
		}
		resp = &types.Response{
			Code: http.StatusOK,
			Data: taskHistory,
		}
		break
	default:
		data := map[string]interface{}{}
		tx := l.svcCtx.Database.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.GetData")).Table(req.Tablename).Where("id = ?", req.Id)
		if role != constant.ROLE_ADMIN {
			tx.Where("tenant IN ?", tenant)
		}
		if err = tx.Take(&data).Error; err != nil {
			l.Logger.Error(err)
			return
		}
		resp = &types.Response{
			Code: http.StatusOK,
			Data: data,
		}
	}
	return
}
