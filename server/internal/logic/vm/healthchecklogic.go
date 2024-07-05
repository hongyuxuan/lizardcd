package vm

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

type HealthcheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHealthcheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthcheckLogic {
	return &HealthcheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HealthcheckLogic) Healthcheck(req *types.HealthCheckReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetTargetAgent(req.Target); err != nil {
		return
	}
	if req.Type == "none" {
		resp = &types.Response{
			Code:    http.StatusOK,
			Message: req.Target + " do not have healthcheck",
		}
	} else {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.VmHealthCheck(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.VmHealthCheck"), &agent.VmHealthCheckRequest{
			Type:   req.Type,
			Method: req.Method,
			Port:   req.Port,
			Uri:    req.Uri,
			Shell:  req.Shell,
		}); err != nil {
			return
		}
		resp = &types.Response{
			Code:    http.StatusOK,
			Data:    string(rpcResponse.Data),
			Message: req.Target + " healthcheck success",
		}
	}
	return
}
