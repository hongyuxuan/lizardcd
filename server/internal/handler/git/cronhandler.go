package git

import (
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/logic/git"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CronHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := git.NewCronLogic(r.Context(), svcCtx)
		resp, err := l.Cron()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
