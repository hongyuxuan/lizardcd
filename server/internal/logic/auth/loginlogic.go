package auth

import (
	"context"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *commontypes.LoginRes, err error) {
	if err = utils.ValidatedUser(req.Username, req.Password, l.svcCtx.Sqlite); err != nil {
		return
	}
	l.Logger.Infof("User \"%s\" login success", req.Username)

	// get user info
	var user commontypes.User
	if err = l.svcCtx.Sqlite.Model(&commontypes.User{}).Where("username = ?", req.Username).WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.GetUser")).First(&user).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	var accessToken string
	var now int64
	if accessToken, now, err = l.svcCtx.GetJwtToken(user, nil); err != nil {
		l.Logger.Error(err)
		return
	}
	resp = &commontypes.LoginRes{
		AccessToken:  accessToken,
		AccessExpire: now + l.svcCtx.Config.Auth.AccessExpire,
		RefreshAfter: now + l.svcCtx.Config.Auth.AccessExpire/2,
	}
	return
}
