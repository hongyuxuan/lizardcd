package istio

import (
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/hongyuxuan/lizardcd/server/internal/logic/istio"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetIstioCrdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResourceReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, err.Error(), nil))
			return
		}

		l := istio.NewGetIstioCrdLogic(r.Context(), svcCtx)
		resp, err := l.GetIstioCrd(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte(resp))
		}
	}
}
