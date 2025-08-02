package task

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RunTaskLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	taskService *svc.TaskService
}

func NewRunTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunTaskLogic {
	return &RunTaskLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		taskService: svc.NewTaskService(ctx, svcCtx),
	}
}

func (l *RunTaskLogic) RunTask(req *commontypes.RunTaskReq) (resp *types.Response, err error) {
	_, _, tenant, _ := utils.GetPayload(l.ctx)
	var application commontypes.Application
	if err = l.svcCtx.Database.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.GetApplication")).
		Model(&commontypes.Application{}).Where("app_name = ?", req.AppName).First(&application).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	taskId, err := l.taskService.RunTask(req, application, tenant[0])
	return &types.Response{
		Code:    http.StatusOK,
		Message: "任务提交成功",
		Data: map[string]string{
			"id": taskId,
		},
	}, err
}
