package lizardcd

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Oauth2Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOauth2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Oauth2Logic {
	return &Oauth2Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Oauth2Logic) Oauth2() (resp *types.Response, err error) {
	var oauth2 []commontypes.Oauth2
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListOauth2")).Find(&oauth2).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: oauth2,
	}
	return
}
