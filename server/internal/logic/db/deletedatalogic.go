package db

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletedataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletedataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletedataLogic {
	return &DeletedataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletedataLogic) Deletedata(req *types.DataByIdReq) (resp *types.Response, err error) {
	data := map[string]interface{}{}
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.DeleteData")).Table(req.Tablename).Where("id = ?", req.Id).Delete(&data).Error; err != nil {
		return
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "删除成功",
	}
	return
}
