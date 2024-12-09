package task

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TaskstatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTaskstatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TaskstatsLogic {
	return &TaskstatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TaskstatsLogic) Taskstats(req *types.TaskStatsReq) (resp *types.Response, err error) {
	res := []map[string]interface{}{}
	if req.By == "month" {
		sub := l.svcCtx.Sqlite.Table("task_history").Select("SUBSTR(init_at, 1, 7) as month").Where("init_at BETWEEN ? and ?", req.TimeFrom, req.TimeTill)
		if err = l.svcCtx.Sqlite.Table("(?) as sub", sub).Select("month, COUNT(*) as count").Group("month").Find(&res).Error; err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		if err = l.svcCtx.Sqlite.Table("task_history").Select(req.By, "COUNT(*) as count").Where("init_at BETWEEN ? and ?", req.TimeFrom, req.TimeTill).Group(req.By).Order("count DESC").Limit(10).Find(&res).Error; err != nil {
			l.Logger.Error(err)
			return
		}
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: res,
	}
	return
}
