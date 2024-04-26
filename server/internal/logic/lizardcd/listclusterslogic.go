package lizardcd

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
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

func (l *ListclustersLogic) Listclusters() (resp *types.Response, err error) {
	res, err := l.listservicesLogic.Listservices()
	var clusterMap = make(map[string][]string)
	for _, service := range res.Data.([]map[string]string) {
		arr := strings.Split(service["service_name"], ".")
		namespace := arr[1]
		cluster := arr[2]
		if _, ok := clusterMap[cluster]; !ok {
			clusterMap[cluster] = []string{}
		}
		if namespace == "*" {
			nss := l.getNamespaces(service["service_name"])
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
