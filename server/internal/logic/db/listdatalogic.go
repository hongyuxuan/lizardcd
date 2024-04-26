package db

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	commonutils "github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListdataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListdataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListdataLogic {
	return &ListdataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListdataLogic) Listdata(req *commontypes.GetDataReq) (resp *types.Response, err error) {
	var count int64
	if req.Tablename == "application" {
		var data []commontypes.Application
		tx := l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, "SpanName", "tidb.ListApplication")).Model(&commontypes.Application{})
		commonutils.SetTx(tx, &count, req)
		if err = tx.Table(req.Tablename).Find(&data).Error; err != nil {
			l.Logger.Error(err)
			return
		}
		resp = &types.Response{
			Code: http.StatusOK,
			Data: map[string]interface{}{
				"total":   count,
				"results": data,
			},
		}
	} else {
		var data []map[string]interface{}
		tx := l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, "SpanName", "sqlite.ListData")).Table(req.Tablename)
		commonutils.SetTx(tx, &count, req)
		if err = tx.Table(req.Tablename).Find(&data).Error; err != nil {
			l.Logger.Error(err)
			return
		}
		if data == nil {
			data = []map[string]interface{}{}
		}
		resp = &types.Response{
			Code: http.StatusOK,
			Data: map[string]interface{}{
				"total":   count,
				"results": data,
			},
		}
	}
	return
}
