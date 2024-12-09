package tekton

import (
	"context"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteTektonLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteTektonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteTektonLogic {
	return &DeleteTektonLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteTektonLogic) DeleteTekton(req *types.ResourceReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	if _, err = ag.DeleteTektonResource(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.ListTektonResource"), &agent.TektonYamlRequest{
		Namespace:    req.Namespace,
		ResourceType: req.ResourceType,
		ResourceName: req.ResourceName,
	}); err != nil {
		l.Logger.Error(err)
		return &types.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, nil
	}
	return &types.Response{
		Code: http.StatusOK,
	}, nil
}
