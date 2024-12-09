package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/oliveagle/jsonpath"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type CallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CallbackLogic {
	return &CallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CallbackLogic) Callback(req *types.CallbackReq) (redirect string, err error) {
	var oauth2 commontypes.Oauth2
	if err = l.svcCtx.Sqlite.Model(&commontypes.Oauth2{}).Where("name = ?", req.Name).First(&oauth2).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	l.Logger.Infof("Oauth2 callback with code \"%s\"", req.Code)
	oauth2Client := utils.NewHttpClient(otel.Tracer("imroc/req"))
	if l.svcCtx.Config.Log.Level == "debug" {
		oauth2Client.EnableDebug(true)
	}
	if strings.HasPrefix(oauth2.AuthorizeUrl, "https") {
		oauth2Client.EnableInsecureSkipVerify()
	}
	if oauth2.HttpProxy != "" {
		oauth2Client.SetProxyURL(oauth2.HttpProxy)
	}

	// 根据code换取access_token
	timestamp := time.Now().UnixNano() / 1e6
	timestampStr := strconv.Itoa(int(timestamp))
	accessRequest := map[string]string{
		"grant_type":      "authorization_code",
		"client_id":       oauth2.ClientId,
		"client_secret":   oauth2.ClientSecret,
		"code":            req.Code,
		"oauth_timestamp": timestampStr,
		"redirect_uri":    oauth2.RedirectUrl,
	}
	var accessToken string
	if strings.Contains(oauth2.TokenUrl, "github.com") { // github返回的access_token不是json格式，需要通过url.Parse解析
		resp := oauth2Client.Post(oauth2.TokenUrl).SetQueryParams(accessRequest).Do(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "http.Oauth2Callback"))
		u, _ := url.Parse("?" + resp.String())
		accessToken = u.Query().Get("access_token")
	} else {
		resp := make(map[string]interface{})
		if err = oauth2Client.Post(oauth2.TokenUrl).SetQueryParams(accessRequest).SetSuccessResult(&resp).Do(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "http.Oauth2Callback")).Err; err != nil {
			l.Logger.Error(err)
			return
		}
		accessToken = resp["access_token"].(string)
	}
	l.Logger.Infof("Get access_token \"%s\" from Oauth2", accessToken)

	// 根据access_token获取profile
	var userinfo interface{}
	if strings.Contains(oauth2.UserinfoUrl, "github.com") { // github Must specify access token via Authorization header
		err = oauth2Client.Get(oauth2.UserinfoUrl).SetBearerAuthToken(accessToken).SetSuccessResult(&userinfo).Do(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "http.GetUserInfo")).Err
	} else {
		err = oauth2Client.Get(oauth2.UserinfoUrl).SetQueryParam("access_token", accessToken).SetSuccessResult(&userinfo).Do(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "http.GetUserInfo")).Err
	}
	if err != nil {
		l.Logger.Error(err)
		return
	}
	userinfob, _ := json.MarshalIndent(userinfo, "", "\t")
	l.Logger.Infof("Get userinfo: \n%s", string(userinfob))
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			err = errorx.NewError(http.StatusUnauthorized, err.Error(), nil)
		}
		return
	}

	// jsonpath提取用户信息
	// 提取username
	var lookupUsername interface{}
	if lookupUsername, err = jsonpath.JsonPathLookup(userinfo, oauth2.UserJsonpath); err != nil {
		l.Logger.Error(err)
		return
	}
	username := utils.AnyToString(lookupUsername)
	// 提取头像
	var lookupAvatar interface{}
	var avatar string
	if oauth2.AvatarJsonpath != "" {
		if lookupAvatar, err = jsonpath.JsonPathLookup(userinfo, oauth2.AvatarJsonpath); err != nil {
			l.Logger.Error(err)
			return
		}
		avatar = utils.AnyToString(lookupAvatar)
	}
	var user commontypes.User
	if err = l.svcCtx.Sqlite.Model(commontypes.User{}).Where("username = ?", username).
		WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "tidb.GetUserByName")).
		First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) { // 用户首次登录
		user = commontypes.User{
			Username: username,
			Role:     "readonly",
			Tenant:   "",
			Profile: map[string]string{
				"avatar": avatar,
			},
			UpdateAt: time.Now(),
		}
		if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "tidb.CreateUser")).Create(&user).Error; err != nil {
			l.Logger.Error(err)
			return
		}
		l.Logger.Infof("Insert new user: %+v into database success", user)
	}
	if accessToken, _, err = l.svcCtx.GetJwtToken(user, nil); err != nil {
		l.Logger.Error(err)
		return
	}
	redirect = fmt.Sprintf("%s?access_token=%s", oauth2.CallbackUrl, accessToken)
	return
}
