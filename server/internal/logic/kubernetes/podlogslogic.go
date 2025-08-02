package kubernetes

import (
	"context"
	"strings"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
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
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	if ag != nil {
		var rpcResponse *agent.YamlResponse
		if rpcResponse, err = ag.GetPodLog(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.GetPodLogs"), &agent.PodLogRequest{
			Namespace:     req.Namespace,
			Podname:       req.PodName,
			ContainerName: req.Container,
			Lines:         uint64(req.Lines),
			Timestamps:    req.Timestamps,
		}); err != nil {
			l.Logger.Error(err)
			return err.Error(), nil
		}
		return strings.TrimSpace(rpcResponse.Data), nil
	} else if ks != nil && ks.IsValid() {
		return ks.GetPodLog(req.Namespace, req.PodName, req.Container, req.Lines, false, req.Timestamps, nil, nil)
	} else {
		return "", errorx.NewDefaultError("Cannot GetLogs of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}
}
