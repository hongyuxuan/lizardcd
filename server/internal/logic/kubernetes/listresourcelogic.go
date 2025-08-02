package kubernetes

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/constant"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commonsvc "github.com/hongyuxuan/lizardcd/common/svc"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	corev1 "k8s.io/api/core/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListResourceLogic {
	return &ListResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListResourceLogic) ListResource(req *types.ListResourceReq) (resp *types.Response, err error) {
	var ag lizardagent.LizardAgent
	var ks *commonsvc.K8sService
	if ag, ks, err = l.svcCtx.GetAgent(req.Cluster, req.Namespace); err != nil {
		return
	}
	resp = &types.Response{
		Code: http.StatusOK,
	}
	var data []byte
	if ag != nil {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.ListResource(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "rpc.ListResource"), &agent.ListResourceRequest{
			Namespace:     req.Namespace,
			ResourceType:  req.ResourceType,
			LabelSelector: req.LabelSelector,
			FieldSelector: req.FieldSelector,
			Continue:      req.Continue,
			Limit:         req.Limit,
		}); err != nil {
			l.Logger.Error(err)
			return
		}
		data = rpcResponse.Data
	} else if ks != nil && ks.IsValid() {
		if data, err = ks.ListResource(req.Namespace, req.ResourceType, req.LabelSelector, req.FieldSelector, req.Continue, req.Limit); err != nil {
			l.Logger.Error(err)
			return
		}
	} else {
		return nil, errorx.NewDefaultError("Cannot ListResource of cluster=%s namespace=%s", req.Cluster, req.Namespace)
	}
	switch req.ResourceType {
	case constant.K8S_RESOURCE_TYPE_PODS:
		r := struct {
			Results  []corev1.Pod `json:"results"`
			Continue string       `json:"continue"`
		}{Results: []corev1.Pod{}}
		json.Unmarshal(data, &r)
		sort.Slice(r.Results, func(i, j int) bool {
			return r.Results[i].CreationTimestamp.After(r.Results[j].CreationTimestamp.Time)
		})
		resp.Data = r
	default:
		var r interface{}
		json.Unmarshal(data, &r)
		resp.Data = r
	}
	return
}
