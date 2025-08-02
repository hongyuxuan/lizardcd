package vm

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SshdeployLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	taskService *svc.TaskService
}

func NewSshdeployLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SshdeployLogic {
	return &SshdeployLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		taskService: svc.NewTaskService(ctx, svcCtx),
	}
}

func (l *SshdeployLogic) Sshdeploy(req *commontypes.SSHDeployReq) (resp *types.Response, err error) {
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "任务提交成功",
	}
	resp.Data, err = l.taskService.SSHDeploy(l.ctx, req)
	return
}
