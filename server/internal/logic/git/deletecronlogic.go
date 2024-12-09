package git

import (
	"context"
	"net/http"
	"strings"

	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletecronLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletecronLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletecronLogic {
	return &DeletecronLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletecronLogic) Deletecron(req *types.SyncStatusReq) (resp *types.Response, err error) {
	for _, appName := range strings.Split(req.Apps, ",") {
		if v, ok := l.svcCtx.CronIdMap[appName]; ok {
			l.svcCtx.Cron.Remove(v.CronId)
		}
		delete(l.svcCtx.CronIdMap, appName)
	}
	return &types.Response{
		Code:    http.StatusOK,
		Message: "删除成功",
	}, nil
}
