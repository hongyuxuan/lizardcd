package db

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatedataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatedataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatedataLogic {
	return &UpdatedataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatedataLogic) Updatedata(req *types.UpdateDataReq) (resp *types.Response, err error) {
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, "SpanName", "tidb.UpdateData")).Table(req.Tablename).Where("id = ?", req.Id).Updates(req.Body).Error; err != nil {
		return
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "更新成功",
	}
	return
}
