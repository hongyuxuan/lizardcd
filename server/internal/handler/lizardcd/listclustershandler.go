package lizardcd

import (
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/logic/lizardcd"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListclustersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := lizardcd.NewListclustersLogic(r.Context(), svcCtx)
		resp, err := l.Listclusters()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
