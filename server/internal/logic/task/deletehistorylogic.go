package task

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteHistoryLogic {
	return &DeleteHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteHistoryLogic) DeleteHistory(req *types.TaskIdReq) (resp *types.Response, err error) {
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.DeleteHistory")).Where("id = ?", req.Id).Delete(&commontypes.TaskHistory{}).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.DeleteHistoryWorkload")).Where("task_history_id = ?", req.Id).Delete(&commontypes.TaskHistoryWorkload{}).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "删除成功",
	}
	return
}
