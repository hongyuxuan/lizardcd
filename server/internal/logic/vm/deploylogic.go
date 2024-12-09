package vm

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeployLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	taskService *svc.TaskService
}

func NewDeployLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeployLogic {
	return &DeployLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		taskService: svc.NewTaskService(ctx, svcCtx),
	}
}

func (l *DeployLogic) Deploy(req *types.VmDeployReq) (resp *types.Response, err error) {
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "任务提交成功",
	}
	resp.Data, err = l.taskService.VmDeploy(l.ctx, req)
	return
}
