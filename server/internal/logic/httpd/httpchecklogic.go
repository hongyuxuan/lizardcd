package httpd

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HttpcheckLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	taskService *svc.TaskService
}

func NewHttpcheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HttpcheckLogic {
	return &HttpcheckLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		taskService: svc.NewTaskService(ctx, svcCtx),
	}
}

func (l *HttpcheckLogic) Httpcheck(req *types.HttpCheckReq) (resp *types.Response, err error) {
	finished, success, message, err := l.taskService.HttpCheck(l.ctx, &req.HttpCheck, req.HttpUrl, req.HttpHeader)
	return &types.Response{
		Code: http.StatusOK,
		Data: types.HttpcheckResponse{
			Finished: finished,
			Success:  success,
			Message:  message,
		},
	}, err
}
