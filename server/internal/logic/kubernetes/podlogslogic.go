package kubernetes

import (
	"context"
	"strings"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PodLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPodLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PodLogsLogic {
	return &PodLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PodLogsLogic) PodLogs(req *types.PodLogReq) (resp string, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.YamlResponse
	if rpcResponse, err = ag.GetPodLog(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetPodLogs"), &agent.PodLogRequest{
		Namespace:     req.Namespace,
		Podname:       req.PodName,
		ContainerName: req.Container,
		Lines:         uint64(req.Lines),
	}); err != nil {
		l.Logger.Error(err)
		return
	}
	resp = strings.TrimSpace(rpcResponse.Data)
	return
}
