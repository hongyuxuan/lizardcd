package helm

import (
	"context"
	"fmt"
	"os"

	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	helmService *commonsvc.HelmService
}

func NewDownloadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadLogic {
	return &DownloadLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		helmService: commonsvc.NewHelmService(ctx),
	}
}

func (l *DownloadLogic) Download(req *types.ShowValuesReq) (file string, err error) {
	destDir := os.TempDir()
	if _, err = l.helmService.Pull(req.RepoUrl, req.ChartName, req.ChartVersion, destDir); err != nil {
		l.Logger.Error(err)
		return
	}
	file = fmt.Sprintf("%s/%s-%s.tgz", destDir, req.ChartName, req.ChartVersion)
	return
}
