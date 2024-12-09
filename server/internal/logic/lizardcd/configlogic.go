package lizardcd

import (
	"context"

	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigLogic {
	return &ConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigLogic) Config() (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
