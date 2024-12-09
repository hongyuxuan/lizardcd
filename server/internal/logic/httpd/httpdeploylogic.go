package httpd

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type HttpdeployLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	taskService *svc.TaskService
}

func NewHttpdeployLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HttpdeployLogic {
	return &HttpdeployLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		taskService: svc.NewTaskService(ctx, svcCtx),
	}
}

func (l *HttpdeployLogic) Httpdeploy(req *types.HttpDeployReq) (resp *types.Response, err error) {
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "任务提交成功",
	}
	resp.Data, err = l.taskService.HttpDeploy(l.ctx, req)
	return
}
