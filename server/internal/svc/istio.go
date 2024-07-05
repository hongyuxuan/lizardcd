package svc

import (
	"context"
	"encoding/json"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/zeromicro/go-zero/core/logx"
	istiometav1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type IstioService struct {
	logx.Logger
	ctx    context.Context
	svcCtx *ServiceContext
}

func NewIstioService(ctx context.Context, svcCtx *ServiceContext) *IstioService {
	return &IstioService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IstioService) SaveDestinationRule(cluster, namespace, appName string, workloads []commontypes.WorkLoad, ag lizardagent.LizardAgent, ifCreate bool) (err error) {
	destinationrule := &v1beta1.DestinationRule{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: namespace,
		},
		Spec: istiometav1.DestinationRule{
			Host:    appName,
			Subsets: []*istiometav1.Subset{},
		},
	}
	for _, workload := range workloads {
		destinationrule.Spec.Subsets = append(destinationrule.Spec.Subsets, &istiometav1.Subset{
			Name: workload.Version,
			Labels: map[string]string{
				"version": workload.Version,
			},
		})
	}
	var rpcResponse *agent.Response
	if ifCreate == true {
		destinationruleBody, _ := json.Marshal(destinationrule)
		if rpcResponse, err = ag.CreateDestinationRule(l.ctx, &lizardagent.IstioCreateRequest{
			Namespace:    namespace,
			ResourceBody: destinationruleBody,
		}); err != nil {
			return
		}
		l.Logger.Infof("Create istio destinationrule success: %+v", string(rpcResponse.Data))
	} else {
		data := map[string]interface{}{
			"spec": destinationrule.Spec,
		}
		body, _ := json.Marshal(data)
		if rpcResponse, err = ag.PatchDestinationRule(l.ctx, &lizardagent.IstioPatchRequest{
			Namespace:    namespace,
			Name:         appName,
			ResourceBody: body,
		}); err != nil {
			return
		}
		l.Logger.Infof("Patch istio destinationrule success: %+v", string(rpcResponse.Data))
	}
	return
}

func (l *IstioService) SaveVirtualService(cluster, namespace, appName, trafficPolicy string, workloads []commontypes.WorkLoad, ag lizardagent.LizardAgent, ifCreate bool) (err error) {
	virtualservice := &v1beta1.VirtualService{
		ObjectMeta: v1.ObjectMeta{
			Name:      appName,
			Namespace: namespace,
		},
		Spec: istiometav1.VirtualService{
			Hosts: []string{appName},
			Http:  []*istiometav1.HTTPRoute{},
		},
	}
	if trafficPolicy == "weight" {
		route := &istiometav1.HTTPRoute{
			Route: []*istiometav1.HTTPRouteDestination{},
		}
		for _, workload := range workloads {
			httpRouteDestination := &istiometav1.HTTPRouteDestination{
				Destination: &istiometav1.Destination{
					Host:   appName,
					Subset: workload.Version,
				},
				Weight: int32(workload.Weight),
			}
			route.Route = append(route.Route, httpRouteDestination)
		}
		virtualservice.Spec.Http = append(virtualservice.Spec.Http, route)
	} else if trafficPolicy == "header" {
		for _, workload := range workloads {
			httpRouteDestination := &istiometav1.HTTPRouteDestination{
				Destination: &istiometav1.Destination{
					Host:   appName,
					Subset: workload.Version,
				},
			}
			route := &istiometav1.HTTPRoute{
				Route: []*istiometav1.HTTPRouteDestination{httpRouteDestination},
			}
			if len(workload.Headers) > 0 {
				route.Match = []*istiometav1.HTTPMatchRequest{}
				for _, header := range workload.Headers {
					h := make(map[string]*istiometav1.StringMatch)
					if header.MatchType == "exact" {
						h[header.Key] = &istiometav1.StringMatch{
							MatchType: &istiometav1.StringMatch_Exact{
								Exact: header.Value,
							},
						}
					} else if header.MatchType == "prefix" {
						h[header.Key] = &istiometav1.StringMatch{
							MatchType: &istiometav1.StringMatch_Prefix{
								Prefix: header.Value,
							},
						}
					} else if header.MatchType == "regex" {
						h[header.Key] = &istiometav1.StringMatch{
							MatchType: &istiometav1.StringMatch_Regex{
								Regex: header.Value,
							},
						}
					}
					route.Match = append(route.Match, &istiometav1.HTTPMatchRequest{
						Headers: h,
					})
				}
			}
			virtualservice.Spec.Http = append(virtualservice.Spec.Http, route)
		}
	}
	var rpcResponse *agent.Response
	if ifCreate == true {
		virtualserviceBody, _ := json.Marshal(virtualservice)
		if rpcResponse, err = ag.CreateVirtualService(l.ctx, &lizardagent.IstioCreateRequest{
			Namespace:    namespace,
			ResourceBody: virtualserviceBody,
		}); err != nil {
			return
		}
		l.Logger.Infof("Create istio virtualservice success: %+v", string(rpcResponse.Data))
	} else {
		data := map[string]interface{}{
			"spec": virtualservice.Spec,
		}
		body, _ := json.Marshal(data)
		if rpcResponse, err = ag.PatchVirtualService(l.ctx, &lizardagent.IstioPatchRequest{
			Namespace:    namespace,
			Name:         appName,
			ResourceBody: body,
		}); err != nil {
			return
		}
		l.Logger.Infof("Patch istio virtualservice success: %+v", string(rpcResponse.Data))
	}
	return
}
