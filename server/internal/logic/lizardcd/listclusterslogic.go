package lizardcd

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListclustersLogic struct {
	logx.Logger
	ctx               context.Context
	svcCtx            *svc.ServiceContext
	listservicesLogic *ListservicesLogic
}

func NewListclustersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListclustersLogic {
	return &ListclustersLogic{
		Logger:            logx.WithContext(ctx),
		ctx:               ctx,
		svcCtx:            svcCtx,
		listservicesLogic: NewListservicesLogic(ctx, svcCtx),
	}
}

func (l *ListclustersLogic) Listclusters(req *types.ListClusterReq) (resp *types.Response, err error) {
	res, err := l.listservicesLogic.Listservices()
	var clusterMap = make(map[string][]string)
	for _, svc := range res.Data.([]map[string]string) {
		meta, err := utils.GetServiceMata(l.svcCtx.Config.ServicePrefix, svc["service_name"])
		if err != nil {
			l.Logger.Error(err)
			continue
		}
		service := meta["Service"]
		cluster := meta["Cluster"]
		namespace := meta["Namespace"]
		if strings.Contains(service, "agent_vm") {
			if !req.WithVm {
				continue
			} else {
				cluster = "vm"
			}
		}
		if _, ok := clusterMap[cluster]; !ok {
			clusterMap[cluster] = []string{}
		}
		if namespace == "*" {
			nss := l.getNamespaces(svc["service_name"])
			clusterMap[cluster] = append(clusterMap[cluster], nss...)
		} else {
			clusterMap[cluster] = append(clusterMap[cluster], namespace)
		}
	}
	resp = &types.Response{
		Code: http.StatusOK,
		Data: clusterMap,
	}
	return
}

func (l *ListclustersLogic) getNamespaces(serviceName string) (res []string) {
	if _, ok := l.svcCtx.AgentList[serviceName]; !ok {
		return
	}
	rpcResponse, err := l.svcCtx.AgentList[serviceName].Client.GetNamespaces(l.ctx, &lizardagent.LabelSelector{LabelSelector: ""})
	if err != nil {
		l.Logger.Error(err)
		return
	} else {
		var r []corev1.Namespace
		json.Unmarshal(rpcResponse.Data, &r)
		for _, ns := range r {
			res = append(res, ns.Name)
		}
	}
	res = lo.Uniq(res)
	return
}
