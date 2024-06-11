package helm

import (
	"context"

	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowValuesLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	helmUtil *utils.HelmUtil
}

func NewShowValuesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowValuesLogic {
	return &ShowValuesLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		helmUtil: utils.NewHelmUtil(ctx),
	}
}

func (l *ShowValuesLogic) ShowValues(req *types.ShowValuesReq) (yaml string, err error) {
	if yaml, err = l.helmUtil.ShowDefaultValues(req.RepoUrl, req.ChartName, req.ChartVersion); err != nil {
		l.Logger.Error(err)
		return
	}
	return
}
