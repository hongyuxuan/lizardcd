package helm

import (
	"context"

	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowReadmeLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	helmUtil *utils.HelmUtil
}

func NewShowReadmeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowReadmeLogic {
	return &ShowReadmeLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		helmUtil: utils.NewHelmUtil(ctx),
	}
}

func (l *ShowReadmeLogic) ShowReadme(req *types.ShowValuesReq) (content string, err error) {
	if content, err = l.helmUtil.ShowReadme(req.RepoUrl, req.ChartName, req.ChartVersion); err != nil {
		l.Logger.Error(err)
		return
	}
	return
}
