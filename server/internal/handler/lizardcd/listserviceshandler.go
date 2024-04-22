package lizardcd

import (
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/logic/lizardcd"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListservicesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := lizardcd.NewListservicesLogic(r.Context(), svcCtx)
		resp, err := l.Listservices()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
