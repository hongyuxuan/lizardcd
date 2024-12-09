package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/constant"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GeneratetokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGeneratetokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GeneratetokenLogic {
	return &GeneratetokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GeneratetokenLogic) Generatetoken(req *types.GenereateTokenReq) (resp *types.Response, err error) {
	_, role, _, _ := utils.GetPayload(l.ctx)
	if role != constant.ROLE_ADMIN {
		return nil, errorx.NewError(http.StatusForbidden, "not permitted", nil)
	}
	var user commontypes.User
	if err = l.svcCtx.Sqlite.First(&user, "id = ?", req.UserId).Error; err != nil {
		return nil, fmt.Errorf("user_id = %d not found", req.UserId)
	}
	var accessToken string
	if accessToken, _, err = l.svcCtx.GetJwtToken(user, &req.ExpireSeconds); err != nil {
		l.Logger.Error(err)
		return
	}
	return &types.Response{
		Code: http.StatusOK,
		Data: accessToken,
	}, nil
}
