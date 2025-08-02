package task

import (
	"context"

	"github.com/golang-module/carbon"
	"github.com/hongyuxuan/lizardcd/common/constant"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type ExecuteTaskLogic struct {
	logx.Logger
	ctx     context.Context
	svcCtx  *svc.ServiceContext
	runTask *RunTaskLogic
}

func NewExecuteTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteTaskLogic {
	return &ExecuteTaskLogic{
		Logger:  logx.WithContext(ctx),
		ctx:     ctx,
		svcCtx:  svcCtx,
		runTask: NewRunTaskLogic(ctx, svcCtx),
	}
}

func (l *ExecuteTaskLogic) ExecuteTask(req *types.ExecuteTaskReq) (resp *types.Response, err error) {
	_, role, tenant, _ := utils.GetPayload(l.ctx)
	var taskHistory commontypes.TaskHistory
	tx := l.svcCtx.Database.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.GetTaskHistory"))
	if role != constant.ROLE_ADMIN {
		tx.Where("tenant IN ?", tenant)
	}
	if err = tx.Preload("TaskHistoryWorkloads").First(&taskHistory, "id = ?", req.Id).Error; err != nil {
		l.Logger.Error(err)
		return

	}
	var artifactUrl string
	if len(taskHistory.TaskHistoryWorkloads) > 0 {
		artifactUrl = taskHistory.TaskHistoryWorkloads[0].Workload.ArtifactUrl
	}
	if req.ArtifactUrl != "" {
		artifactUrl = req.ArtifactUrl
	}
	runTaskReq := &commontypes.RunTaskReq{
		Id:          taskHistory.Id,
		AppName:     taskHistory.AppName,
		TaskType:    taskHistory.TaskType,
		TriggerType: taskHistory.TriggerType,
		Labels:      taskHistory.Labels,
		Workloads: lo.Map(taskHistory.TaskHistoryWorkloads, func(w commontypes.TaskHistoryWorkload, _ int) commontypes.Workload {
			w.Workload.ArtifactUrl = artifactUrl
			return w.Workload
		}),
		InitAt:      carbon.FromStdTime(taskHistory.InitAt.Time).Format("Y-m-d H:i:s"),
		ArtifactUrl: artifactUrl,
		Waiting:     false,
	}
	return l.runTask.RunTask(runTaskReq)
}
