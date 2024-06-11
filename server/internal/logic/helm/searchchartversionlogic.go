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

type SearchChartVersionLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	helmUtil *utils.HelmUtil
}

func NewSearchChartVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchChartVersionLogic {
	return &SearchChartVersionLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		helmUtil: utils.NewHelmUtil(ctx),
	}
}

func (l *SearchChartVersionLogic) SearchChartVersion(req *types.ChartReq) (resp *types.Response, err error) {
	var res []*commontypes.ChartListResponse
	if res, err = l.helmUtil.SearchChartVersions(req.Name, req.ChartName); err != nil {
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: res,
	}
	return
}
