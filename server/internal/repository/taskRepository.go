package repository

import (
	"context"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type TaskRepository struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskRepository(ctx context.Context, svcCtx *svc.ServiceContext) *TaskRepository {
	return &TaskRepository{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (t *TaskRepository) DeleteTaskHistory(taskHistory commontypes.TaskHistory) (err error) {
	if err = t.svcCtx.Database.WithContext(context.WithValue(t.ctx, commontypes.TraceIDKey{}, "sqlite.DeleteHistoryWorkload")).Delete(&commontypes.TaskHistoryWorkload{}, "task_history_id = ?", taskHistory.Id).Error; err != nil {
		t.Logger.Error(err)
		return
	}
	t.Logger.Infof("Successfully delete task_history_workload for task_history_id=%s", taskHistory.Id)
	if err = t.svcCtx.Database.WithContext(context.WithValue(t.ctx, commontypes.TraceIDKey{}, "sqlite.DeleteTaskHistory")).Delete(&taskHistory).Error; err != nil {
		t.Logger.Error(err)
		return
	}
	t.Logger.Infof("Successfully delete task_history for id=%s", taskHistory.Id)
	return
}
