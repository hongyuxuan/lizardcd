package tekton

import (
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/logic/tekton"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ApplyTektonHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PatchVariableReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, err.Error(), nil))
			return
		}

		l := tekton.NewApplyTektonLogic(r.Context(), svcCtx)
		resp, err := l.ApplyTekton(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
