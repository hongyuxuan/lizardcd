package lizardcd

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListtargetsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListtargetsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListtargetsLogic {
	return &ListtargetsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListtargetsLogic) Listtargets() (resp *types.Response, err error) {
	_, role, _, namespaces := utils.GetPayload(l.ctx)
	targets := []map[string]interface{}{}
	for k, v := range l.svcCtx.AgentList {
		target, err := utils.GetTarget(l.svcCtx.Config.ServicePrefix, k, namespaces, role)
		if err != nil {
			continue
		}
		targets = append(targets, map[string]interface{}{
			"ip":     target,
			"labels": v.Labels,
		})
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: targets,
	}
	return
}
