package db

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
