package db

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetdataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetdataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetdataLogic {
	return &GetdataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetdataLogic) Getdata(req *types.DataByIdReq) (resp *types.Response, err error) {
	_, _, tenant, _ := utils.GetPayload(l.ctx)
	if req.Tablename == "task_history" {
		var taskHistory commontypes.TaskHistory
		if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, "SpanName", "sqlite.GetTaskHistory")).Where("tenant = ?", tenant).Preload("TaskHistoryWorkloads").First(&taskHistory, "id = ?", req.Id).Error; err != nil {
			l.Logger.Error(err)
			return
		}
		resp = &types.Response{
			Code: http.StatusOK,
			Data: taskHistory,
		}
	}
	return
}
