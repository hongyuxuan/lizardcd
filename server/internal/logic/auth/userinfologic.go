package auth

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserinfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserinfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserinfoLogic {
	return &UserinfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserinfoLogic) Userinfo() (resp *types.Response, err error) {
	payloads := l.ctx.Value("payloads").(map[string]interface{})
	var user commontypes.User
	if err = l.svcCtx.Sqlite.Model(&commontypes.User{}).Where("username = ?", payloads["username"]).First(&user).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: user,
	}
	return
}
