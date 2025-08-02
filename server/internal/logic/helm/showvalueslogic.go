package helm

import (
	"context"

	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShowValuesLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	helmService *commonsvc.HelmService
}

func NewShowValuesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShowValuesLogic {
	return &ShowValuesLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		helmService: commonsvc.NewHelmService(ctx),
	}
}

func (l *ShowValuesLogic) ShowValues(req *types.ShowValuesReq) (yaml string, err error) {
	if yaml, err = l.helmService.ShowDefaultValues(req.RepoUrl, req.ChartName, req.ChartVersion); err != nil {
		l.Logger.Error(err)
		return
	}
	return
}
