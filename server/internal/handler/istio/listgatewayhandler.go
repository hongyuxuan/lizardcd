package istio

import (
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/logic/istio"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListGatewayHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListWorkloadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := istio.NewListGatewayLogic(r.Context(), svcCtx)
		resp, err := l.ListGateway(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
