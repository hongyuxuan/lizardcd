package lizardcd

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/zeromicro/go-zero/core/logx"

	corev1 "k8s.io/api/core/v1"
)

type PodEventsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPodEventsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PodEventsLogic {
	return &PodEventsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PodEventsLogic) PodEvents(req *types.GetPodEventReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	var rpcResponse *agent.Response
	if rpcResponse, err = ag.GetPodEvent(context.WithValue(l.ctx, "SpanName", "rpc.PodEvents"), &agent.GetPodEventRequest{
		Namespace: req.Namespace,
		PodName:   req.PodName,
	}); err != nil {
		l.Logger.Error(err)
		return
	}
	var r []corev1.Event
	json.Unmarshal(rpcResponse.Data, &r)
	resp = &types.Response{
		Code: http.StatusOK,
		Data: r,
	}
	return
}
