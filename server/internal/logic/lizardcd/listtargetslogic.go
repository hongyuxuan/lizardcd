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
	var targets []string
	for k := range l.svcCtx.AgentList {
		target, err := utils.GetTarget(l.svcCtx.Config.ServicePrefix, k, namespaces, role)
		if err != nil {
			continue
		}
		targets = append(targets, target)
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: targets,
	}
	return
}
