package helm

import (
	"context"
	"net/http"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchChartLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	helmUtil *utils.HelmUtil
}

func NewSearchChartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchChartLogic {
	return &SearchChartLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		helmUtil: utils.NewHelmUtil(ctx),
	}
}

func (l *SearchChartLogic) SearchChart(req *types.RepoReq) (resp *types.Response, err error) {
	var entry *commontypes.HelmRepositories
	if err = l.svcCtx.Sqlite.Where("name = ?", req.RepoName).First(&entry).Error; err != nil {
		return
	}
	var res []*commontypes.ChartListResponse
	if res, err = l.helmUtil.SearchChart(entry.URL, req.RepoName, req.ChartName); err != nil {
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: res,
	}
	return
}
