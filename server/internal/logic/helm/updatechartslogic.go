package helm

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"helm.sh/helm/v3/pkg/repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateChartsLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	helmUtil *utils.HelmUtil
}

func NewUpdateChartsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateChartsLogic {
	return &UpdateChartsLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		helmUtil: utils.NewHelmUtil(ctx),
	}
}

func (l *UpdateChartsLogic) UpdateCharts(req *types.ListWorkloadReq) (resp *types.Response, err error) {
	_, _, tenant, _ := utils.GetPayload(l.ctx)
	var res []*commontypes.HelmRepositories
	if err = l.svcCtx.Sqlite.Where("tenant = ?", tenant).Find(&res).Error; err != nil {
		l.Logger.Error(err)
		return
	}
	var entries []*repo.Entry
	for _, r := range res {
		entries = append(entries, &repo.Entry{
			Name:     r.Name,
			URL:      r.URL,
			Username: r.Username,
			Password: r.Password,
		})
	}
	// update repo for server
	if err = l.helmUtil.Update(entries); err != nil {
		l.Logger.Error(err)
		return
	}

	// update repo for agent
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	b, _ := json.Marshal(res)
	if _, err = ag.HelmUpdateRepo(l.ctx, &lizardagent.HelmEntriesRequest{Entries: b}); err != nil {
		l.Logger.Error(err)
		return
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "更新成功",
	}
	return
}
