package helm

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/constant"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRepoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRepoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRepoLogic {
	return &ListRepoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListRepoLogic) ListRepo() (resp *types.Response, err error) {
	_, role, tenant, _ := utils.GetPayload(l.ctx)
	var res []*commontypes.HelmRepositories
	if role != constant.ROLE_ADMIN {
		if err = l.svcCtx.Sqlite.Where("tenant = ?", tenant).Find(&res).Error; err != nil {
			l.Logger.Error(err)
			return
		}
	} else if err = l.svcCtx.Sqlite.Find(&res).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: res,
	}
	return
}
