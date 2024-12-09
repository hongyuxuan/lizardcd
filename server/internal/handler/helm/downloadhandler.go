package helm

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/hongyuxuan/lizardcd/server/internal/logic/helm"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DownloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ShowValuesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, err.Error(), nil))
			return
		}

		l := helm.NewDownloadLogic(r.Context(), svcCtx)
		file, err := l.Download(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			w.Header().Add("Content-Disposition", "attachment; filename="+filepath.Base(file))
			http.ServeFile(w, r, file)
			os.Remove(file)
		}
	}
}
