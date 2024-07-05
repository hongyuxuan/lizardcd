package static

import (
	"context"

	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DocfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDocfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocfileLogic {
	return &DocfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocfileLogic) Docfile(req *types.StaticReq) error {
	// todo: add your logic here and delete this line

	return nil
}
