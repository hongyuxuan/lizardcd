package db

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/repository"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletedataLogic struct {
	logx.Logger
	ctx                   context.Context
	svcCtx                *svc.ServiceContext
	applicationRepository *repository.ApplicationRepository
	taskRepository        *repository.TaskRepository
}

func NewDeletedataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletedataLogic {
	return &DeletedataLogic{
		Logger:                logx.WithContext(ctx),
		ctx:                   ctx,
		svcCtx:                svcCtx,
		applicationRepository: repository.NewApplicationRepository(ctx, svcCtx),
		taskRepository:        repository.NewTaskRepository(ctx, svcCtx),
	}
}

func (l *DeletedataLogic) Deletedata(req *types.DataByIdReq) (resp *types.Response, err error) {
	_, _, tenant, _ := utils.GetPayload(l.ctx)
	if req.Tablename == "application" {
		if err = l.applicationRepository.Delete(req.Id, tenant[0]); err != nil {
			l.Logger.Error(err)
			return
		}
	} else if req.Tablename == "task_history" {
		if err = l.taskRepository.DeleteTaskHistory(commontypes.TaskHistory{Id: req.Id}); err != nil {
			return
		}
		return &types.Response{
			Code: http.StatusOK,
		}, nil
	} else {
		data := map[string]interface{}{}
		if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.DeleteData")).
			Table(req.Tablename).
			Where("id = ?", req.Id).
			Delete(&data).Error; err != nil {
			return
		}
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "删除成功",
	}
	return
}
