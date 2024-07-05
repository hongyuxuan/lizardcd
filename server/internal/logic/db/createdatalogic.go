package db

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatedataLogic struct {
	logx.Logger
	ctx          context.Context
	svcCtx       *svc.ServiceContext
	istioService *svc.IstioService
}

func NewCreatedataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatedataLogic {
	return &CreatedataLogic{
		Logger:       logx.WithContext(ctx),
		ctx:          ctx,
		svcCtx:       svcCtx,
		istioService: svc.NewIstioService(ctx, svcCtx),
	}
}

func (l *CreatedataLogic) Createdata(req *types.CreateDataReq) (resp *types.Response, err error) {
	if err = l.svcCtx.Sqlite.WithContext(context.WithValue(l.ctx, commontypes.TraceIDKey{}, "sqlite.CreateData")).Table(req.Tablename).Create(&req.Body).Error; err != nil {
		return
	}
	// create application may also create a istio CRD
	if req.Tablename == "application" && req.Body["enable_traffic_control"].(bool) == true {
		workloadStr := req.Body["workload"].(string)
		var workloads []commontypes.WorkLoad
		json.Unmarshal([]byte(workloadStr), &workloads)
		cluster := workloads[0].Cluster
		namespace := workloads[0].Namespace
		appName := req.Body["app_name"].(string)
		var ag lizardagent.LizardAgent
		if ag, err = l.svcCtx.GetAgent(cluster, namespace); err != nil {
			return
		}
		if err = l.istioService.SaveDestinationRule(cluster, namespace, appName, workloads, ag, true); err != nil { // create destinationrule
			l.Logger.Error(err)
		} else {
			trafficPolicy := req.Body["traffic_policy"].(string)
			if err = l.istioService.SaveVirtualService(cluster, namespace, appName, trafficPolicy, workloads, ag, true); err != nil { // create virtualservice
				l.Logger.Error(err)
			}
		}
	}
	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "新增成功",
	}
	return
}
