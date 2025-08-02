package task

import (
	"context"
	"net/http"
	"strings"
	"sync"

	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApphistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApphistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApphistoryLogic {
	return &ApphistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApphistoryLogic) Apphistory(req *types.AppHistoryReq) (resp *types.Response, err error) {
	_, role, tenant, _ := utils.GetPayload(l.ctx)
	data := make(map[string][]commontypes.TaskHistory)
	if req.Apps != "" {
		apps := strings.Split(req.Apps, ",")
		var wg sync.WaitGroup
		resultChan := make(chan []commontypes.TaskHistory, len(apps))
		for _, app := range apps {
			wg.Add(1)
			go func(app string, size int, sort string, ch chan []commontypes.TaskHistory) {
				var history []commontypes.TaskHistory
				tx := l.svcCtx.Database.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.ListTaskHistory")).Model(commontypes.TaskHistory{})
				var count int64
				utils.SetTx(tx, &commontypes.GetDataReq{
					Tablename: "task_history",
					Page:      1,
					Size:      size,
					Filter:    "app_name==" + app,
					Sort:      sort,
				}, &count, role, tenant, nil)
				if err = tx.Find(&history).Error; err != nil {
					l.Logger.Error(err)
				}
				wg.Done()
				ch <- history
			}(app, req.Size, req.Sort, resultChan)
		}
		wg.Wait()
		close(resultChan)
		for result := range resultChan {
			if len(result) > 0 {
				appName := result[0].AppName
				data[appName] = result
			}
		}
	}
	return &types.Response{
		Code: http.StatusOK,
		Data: data,
	}, nil
}
