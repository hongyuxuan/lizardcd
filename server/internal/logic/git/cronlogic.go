package git

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/robfig/cron/v3"

	"github.com/zeromicro/go-zero/core/logx"
)

type CronLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCronLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CronLogic {
	return &CronLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CronLogic) Cron() (resp *types.Response, err error) {
	var cronIds []cron.EntryID
	for _, entry := range l.svcCtx.Cron.Entries() {
		cronIds = append(cronIds, entry.ID)
	}
	return &types.Response{
		Code: http.StatusOK,
		Data: map[string]interface{}{
			"cronId:":     cronIds,
			"application": l.svcCtx.CronIdMap,
		},
	}, nil
}
