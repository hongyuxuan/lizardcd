package task

import (
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/hongyuxuan/lizardcd/server/internal/logic/task"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RunTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RunTaskReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, err.Error(), nil))
			return
		}

		l := task.NewRunTaskLogic(r.Context(), svcCtx)
		resp, err := l.RunTask(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
