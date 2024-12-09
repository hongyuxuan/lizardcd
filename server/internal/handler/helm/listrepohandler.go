package helm

import (
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/hongyuxuan/lizardcd/server/internal/logic/helm"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListRepoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := helm.NewListRepoLogic(r.Context(), svcCtx)
		resp, err := l.ListRepo()
		if err != nil {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, err.Error(), nil))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
