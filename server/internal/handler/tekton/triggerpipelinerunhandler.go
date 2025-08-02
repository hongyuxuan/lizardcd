package tekton

import (
	"net/http"

	"github.com/hongyuxuan/lizardcd/server/internal/logic/tekton"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func TriggerpipelinerunHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TektontriggerReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		l := tekton.NewTriggerpipelinerunLogic(r.Context(), svcCtx)
		resp, err := l.Triggerpipelinerun(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
