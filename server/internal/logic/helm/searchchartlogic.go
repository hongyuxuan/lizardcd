package helm

import (
	"context"
	"net/http"

	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchChartLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	helmService *commonsvc.HelmService
}

func NewSearchChartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchChartLogic {
	return &SearchChartLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		helmService: commonsvc.NewHelmService(ctx),
	}
}

func (l *SearchChartLogic) SearchChart(req *types.RepoReq) (resp *types.Response, err error) {
	var entry *commontypes.HelmRepositories
	if err = l.svcCtx.Database.Where("name = ?", req.RepoName).First(&entry).Error; err != nil {
		return
	}
	var res []*commontypes.ChartListResponse
	if res, err = l.helmService.SearchChart(entry.URL, req.RepoName, req.ChartName); err != nil {
		return
	}
	if res == nil {
		res = []*commontypes.ChartListResponse{}
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: res,
	}
	return
}
