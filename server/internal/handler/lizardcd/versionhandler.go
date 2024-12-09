package lizardcd

import (
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/svc"
)

func VersionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(svcCtx.Version))
	}
}
