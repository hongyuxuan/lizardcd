package svc

import (
	"context"
	"strings"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/zeromicro/go-zero/core/logx"
)

type WsOperator struct {
	logx.Logger
	svcCtx *ServiceContext
	cancel context.CancelFunc
}

func NewWsOperator(ctx context.Context, svcCtx *ServiceContext) *WsOperator {
	return &WsOperator{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (ws *WsOperator) GetPodLogs(msg types.WsMessage) {
	cluster := msg.MessageData["cluster"].(string)
	namespace := msg.MessageData["namespace"].(string)
	podname := msg.MessageData["pod_name"].(string)
	container := msg.MessageData["container"].(string)
	lines := msg.MessageData["lines"].(float64)
	follow := msg.MessageData["follow"].(bool)
	var ag lizardagent.LizardAgent
	var err error
	if ag, err = ws.svcCtx.GetAgent(cluster, namespace); err != nil {
		return
	}
	if follow {
		var ctx context.Context
		ctx, ws.cancel = context.WithCancel(context.Background())
		stream, e := ag.GetPodLogFollow(ctx, &agent.PodLogRequest{
			Namespace:     namespace,
			Podname:       podname,
			ContainerName: container,
			Lines:         uint64(lines),
		})
		if e != nil {
			ws.Logger.Error(err)
			return
		}
		go ws.recvPodLogs(stream, msg)
	} else if ws.cancel != nil {
		ws.cancel()
	}
}

func (ws *WsOperator) recvPodLogs(stream agent.LizardAgent_GetPodLogFollowClient, msg types.WsMessage) {
	for {
		res, e := stream.Recv()
		if e != nil && strings.Contains(e.Error(), "EOF") {
			ws.Logger.Infof("Lizardcd agent closed")
			break
		}
		if e != nil && strings.Contains(e.Error(), "context canceled") {
			ws.Logger.Infof("Context canceled and stop recv pod logs")
			break
		}
		if e != nil {
			ws.Logger.Errorf("Recv pod log error: %v", e.Error())
			break
		}
		ws.Logger.Debug("----->", res.Data)
		ws.svcCtx.Hub.broadcast <- types.WsMessage{
			MessageType: msg.MessageType,
			MessageTo:   msg.MessageTo,
			MessageData: map[string]interface{}{
				"content": res.Data,
			},
		}
	}
}
