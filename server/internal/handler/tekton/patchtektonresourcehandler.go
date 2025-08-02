package tekton

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/hongyuxuan/lizardcd/server/internal/logic/tekton"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	tektontypes "github.com/hongyuxuan/tekton-sdk-go/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func PatchTektonResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`^/lizardcd/tekton/cluster/([^/]+)/namespace/([^/]+)/([^/]+)/([^/]+)$`)
		matches := re.FindStringSubmatch(r.URL.Path)
		req := types.ResourceReq{}
		if len(matches) == 5 {
			req.Cluster = matches[1]
			req.Namespace = matches[2]
			req.ResourceType = matches[3]
			req.ResourceName = matches[4]
		} else {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, "path params error", nil))
			return
		}

		var data []tektontypes.PatchOptions
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields() // 不允许未知字段
		if err := decoder.Decode(&data); err != nil {
			httpx.Error(w, errorx.NewError(http.StatusBadRequest, err.Error(), nil))
			return
		}

		l := tekton.NewPatchTektonResourceLogic(r.Context(), svcCtx)
		resp, err := l.PatchTektonResource(&req, data)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
