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

type SearchChartVersionLogic struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *svc.ServiceContext
	helmService *commonsvc.HelmService
}

func NewSearchChartVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchChartVersionLogic {
	return &SearchChartVersionLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		helmService: commonsvc.NewHelmService(ctx),
	}
}

func (l *SearchChartVersionLogic) SearchChartVersion(req *types.ChartReq) (resp *types.Response, err error) {
	var res []*commontypes.ChartListResponse
	if res, err = l.helmService.SearchChartVersions(req.Name, req.ChartName); err != nil {
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: res,
	}
	return
}
