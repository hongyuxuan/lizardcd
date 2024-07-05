package auth

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdduserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdduserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdduserLogic {
	return &AdduserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdduserLogic) Adduser(req *types.AddUserReq) (resp *types.Response, err error) {
	generatedPassword := utils.GenerateRandomString(10)
	if err = utils.AddUser(req.Username, generatedPassword, req.Role, req.Tenant, l.svcCtx.Sqlite); err != nil {
		err = errorx.NewDefaultError("Failed to create user \"%s\": %v", req.Username, err)
		l.Logger.Error(err)
		return
	} else {
		l.Logger.Infof("Successfully create user \"%s\" with password \"\"", req.Username, generatedPassword)
	}
	// add user must also add related settings
	utils.AddSettings(req.Tenant, l.svcCtx.Sqlite)

	resp = &types.Response{
		Code:    http.StatusOK,
		Data:    generatedPassword,
		Message: "新建用户成功",
	}
	return
}
