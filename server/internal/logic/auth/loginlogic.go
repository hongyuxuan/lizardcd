package auth

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/samber/lo"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type NsCluster struct {
	Namespace string `json:"namespace"`
	Cluster   string `json:"cluster"`
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

	// get user info
	var user struct {
		commontypes.User
		Namespaces string
	}
	if err = l.svcCtx.Sqlite.Model(&commontypes.User{}).Select("user.*,tenant.namespaces").Joins("left join tenant on user.tenant = tenant.tenant_name").Where("username = ?", req.Username).First(&user).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	var nscluster []NsCluster
	if err = json.Unmarshal([]byte(user.Namespaces), &nscluster); err != nil {
		l.Logger.Error(err)
		return
	}
	ns := lo.Map(nscluster, func(item NsCluster, _ int) string {
		return item.Namespace
	})
	payloads := map[string]interface{}{
		"username":  req.Username,
		"role":      user.Role,
		"tenant":    user.Tenant,
		"namespace": strings.Join(ns, ","),
	}

	var accessExpire = l.svcCtx.Config.Auth.AccessExpire
	var accessToken string
	now := time.Now().Unix()
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
