package db

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatedataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatedataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatedataLogic {
	return &CreatedataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatedataLogic) Createdata(req *types.CreateDataReq) (resp *types.Response, err error) {
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, "SpanName", "sqlite.CreateData")).Table(req.Tablename).Create(&req.Body).Error; err != nil {
		return
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "新增成功",
	}
	return
}
