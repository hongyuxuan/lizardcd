package helm

import (
	"context"
	"fmt"

	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	helmUtil *utils.HelmUtil
}

func NewDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadLogic {
	return &DownloadLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		helmUtil: utils.NewHelmUtil(ctx),
	}
}

func (l *DownloadLogic) Download(req *types.ShowValuesReq) (file string, err error) {
	destDir := "./cache"
	if _, err = l.helmUtil.Pull(req.RepoUrl, req.ChartName, req.ChartVersion, destDir); err != nil {
		l.Logger.Error(err)
		return
	}
	file = fmt.Sprintf("%s/%s-%s.tgz", destDir, req.ChartName, req.ChartVersion)
	return
}
