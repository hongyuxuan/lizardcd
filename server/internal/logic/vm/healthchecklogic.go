package vm

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthcheckLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	taskService *svc.TaskService
}

func NewHealthcheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthcheckLogic {
	return &HealthcheckLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		taskService: svc.NewTaskService(ctx, svcCtx),
	}
}

func (l *HealthcheckLogic) Healthcheck(req *types.HealthCheckReq) (resp *types.Response, err error) {
	resp = &types.Response{
		Code: http.StatusOK,
	}
	resp.Message, resp.Data, err = l.taskService.Healthcheck(l.ctx, &req.HealthCheck, req.Target)
	return
}
