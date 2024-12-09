package tekton

import (
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/hongyuxuan/lizardcd/server/internal/logic/tekton"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TriggerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TektontriggerReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, err.Error(), nil))
			return
		}
		secret := r.Header.Get("X-Gitlab-Token")
		if secret == "" {
			httpx.Error(w, errorx.NewDefaultError("HTTP header X-Gitlab-Token cannot be empty"))
			return
		}
		l := tekton.NewTriggerLogic(r.Context(), svcCtx)
		resp, err := l.Trigger(&req, secret)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
