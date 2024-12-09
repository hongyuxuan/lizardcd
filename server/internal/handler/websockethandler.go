package handler

import (
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/svc"
)

func websocketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		svc.ServeWs(svcCtx, w, r)
	}
}
