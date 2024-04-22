package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
	l.Logger.Infof("user \"%s\" login success", req.Username)
	var accessExpire = l.svcCtx.Config.Auth.AccessExpire
	var accessToken string
	now := time.Now().Unix()
	payloads := map[string]interface{}{
		"username": req.Username,
	}
	if accessToken, err = l.getToken(now, payloads, accessExpire); err != nil {
		return nil, err
	}
	resp = &commontypes.LoginRes{
		AccessToken:  accessToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}
	return
}

func (l *LoginLogic) getToken(iat int64, payloads map[string]interface{}, seconds int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["payloads"] = payloads
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))
}
