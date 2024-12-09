package lizardcd

import (
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/logic/lizardcd"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func Oauth2Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := lizardcd.NewOauth2Logic(r.Context(), svcCtx)
		resp, err := l.Oauth2()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
