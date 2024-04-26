package auth

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChpasswdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChpasswdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChpasswdLogic {
	return &ChpasswdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChpasswdLogic) Chpasswd(req *types.ChpasswdReq) (resp *types.Response, err error) {
	if err = utils.ModifyPassword(req.Username, req.OldPassword, req.NewPassword, l.svcCtx.Sqlite); err != nil {
		l.Logger.Error(err)
		return
	}
	l.Logger.Infof("user \"%s\" changed password success")
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "修改密码成功",
	}
	return
}
