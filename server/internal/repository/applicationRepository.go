package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/common/constant"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/logic/kubernetes"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
)

type ApplicationRepository struct {
	logx.Logger
	ctx             context.Context
	svcCtx          *svc.ServiceContext
	application     *commontypes.Application
	applicationFaas *commontypes.ApplicationFaas
	istioService    *svc.IstioService
	gitService      *svc.GitService
	taskService     *svc.TaskService
	cronService     *svc.CronService
	variableLogic   *kubernetes.PatchVariableLogic
	taskRepository  *TaskRepository
}

func NewApplicationRepository(ctx context.Context, svcCtx *svc.ServiceContext) *ApplicationRepository {
	return &ApplicationRepository{
		Logger:          logx.WithContext(ctx),
		ctx:             ctx,
		svcCtx:          svcCtx,
		application:     &commontypes.Application{},
		applicationFaas: &commontypes.ApplicationFaas{},
		istioService:    svc.NewIstioService(ctx, svcCtx),
		gitService:      svc.NewGitService(ctx, svcCtx),
		taskService:     svc.NewTaskService(ctx, svcCtx),
		cronService:     svc.NewCronService(ctx, svcCtx),
		variableLogic:   kubernetes.NewPatchVariableLogic(ctx, svcCtx),
		taskRepository:  NewTaskRepository(ctx, svcCtx),
	}
}

func (r *ApplicationRepository) Save(body map[string]interface{}, tenant string, ifCreate bool) (id interface{}, err error) {
	b, _ := json.Marshal(body)
	if err = json.Unmarshal(b, &r.application); err != nil {
		r.Logger.Error(err)
		return
	}
	if err = r.svcCtx.Database.WithContext(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "sqlite.SaveApplication")).
		Save(&r.application).Error; err != nil {
		r.Logger.Error(err)
		return
	}
	if r.application.EnableTrafficControl {
		r.saveIstio(r.application.AppName, r.application.TrafficPolicy, r.application.Workload, ifCreate)
	}
	r.cronService.RemoveCron(*r.application, tenant, false)
	if r.application.GitOps != nil {
		// parse yaml to k8s object
		var yamlstring string
		if _, yamlstring, err = r.gitService.GetYamlFromGit(*r.application, tenant); err != nil {
			r.Logger.Error(err)
			return nil, errors.Unwrap(err)
		}
		if err = r.gitService.ParseYaml(r.ctx, *r.application, yamlstring); err != nil {
			r.Logger.Error(err)
			return nil, errors.Unwrap(err)
		}
		// add autosync
		if r.application.GitOps.SyncType == constant.APP_SYNC_TYPE_AUTO {
			if err = r.cronService.AddGitOpsCron(*r.application, tenant); err != nil {
				r.Logger.Error(err)
				return
			}
		}
	}
	return r.application.Id, nil
}

func (r *ApplicationRepository) SaveFaas(body map[string]interface{}, ifCreate bool) (id interface{}, err error) {
	b, _ := json.Marshal(body)
	json.Unmarshal(b, &r.applicationFaas)
	if err = r.svcCtx.Database.WithContext(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "sqlite.CreateApplicationFaas")).
		Save(&r.applicationFaas).Error; err != nil {
		return
	}
	// generate workloads from faas
	b, _ = json.Marshal(r.applicationFaas.Variables["Versions"])
	var faasVersion []commontypes.FaasVersion
	json.Unmarshal(b, &faasVersion)
	appName := r.applicationFaas.Variables["Appname"].(string)
	workloads := lo.Map(faasVersion, func(item commontypes.FaasVersion, _ int) commontypes.Workload {
		return commontypes.Workload{
			Cluster:       r.applicationFaas.Cluster,
			Namespace:     r.applicationFaas.Namespace,
			WorkloadType:  constant.K8S_RESOURCE_TYPE_DEPLOYMENTS,
			WorkloadName:  appName + "-" + item.Version,
			ContainerName: appName + "-container",
			Version:       item.Version,
			Weight:        item.Weight,
			Enable:        true,
		}
	})
	if len(workloads) == 0 {
		workloads = []commontypes.Workload{
			{
				Cluster:       r.applicationFaas.Cluster,
				Namespace:     r.applicationFaas.Namespace,
				WorkloadType:  constant.K8S_RESOURCE_TYPE_DEPLOYMENTS,
				WorkloadName:  appName,
				ContainerName: appName + "-container",
				Enable:        true,
			},
		}
	}
	r.svcCtx.Database.WithContext(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "sqlite.UpdateApplication")).
		Table("application").
		Where("id = ?", r.applicationFaas.ApplicationId).
		Updates(commontypes.Application{
			Workload: workloads,
		})
	r.Logger.Infof("update application id=%d workloads=%+v success", r.applicationFaas.ApplicationId, workloads)
	// create k8s resource
	r.variableLogic.PatchVariable(&types.PatchVariableReq{
		Cluster:   r.applicationFaas.Cluster,
		Namespace: r.applicationFaas.Namespace,
		Content:   r.applicationFaas.Template,
		Variables: r.applicationFaas.Variables,
	})
	// // create istio resource
	if len(faasVersion) > 0 {
		r.saveIstio(appName, constant.APP_TRAFFIC_POLICY_WEIGHT, workloads, ifCreate)
	}
	return r.applicationFaas.Id, nil
}

func (r *ApplicationRepository) Delete(id, tenant string) (err error) {
	var application commontypes.Application
	if err = r.svcCtx.Database.WithContext(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "sqlite.GetApplication")).First(&application, "id = ?", id).Error; err != nil {
		r.Logger.Error(err)
		return
	}
	// delete application
	if err = r.svcCtx.Database.WithContext(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "sqlite.DeleteApplication")).Delete(&commontypes.Application{}, id).Error; err != nil {
		return
	}
	r.Logger.Infof("Successfully delete application %s", application.AppName)
	// delete autosync from etcd
	r.cronService.RemoveCron(application, tenant, true)
	// prune k8s resource
	if application.GitOps != nil && application.GitOps.PruneOnDelete {
		var yamlstring string
		if _, yamlstring, err = r.gitService.GetYamlFromGit(application, tenant); err != nil {
			r.Logger.Errorf("Error in delete application GetYamlFromGit: %w", err)
			return errors.Unwrap(err)
		}
		var ag lizardagent.LizardAgent
		for _, workload := range application.Workload {
			if ag, _, err = r.svcCtx.GetAgent(workload.Cluster, workload.Namespace); err != nil {
				return
			}
			if _, err = ag.DeleteYaml(r.ctx, &lizardagent.YamlRequest{
				Namespace: workload.Namespace,
				Ymlstring: yamlstring,
			}); err != nil {
				r.Logger.Errorf("Error in delete yaml: %w from cluster=%s namespace=%s", err, workload.Cluster, workload.Namespace)
			} else {
				r.Logger.Infof("Successfully delete yaml for application=%s from cluster=%s namespace=%s", application.AppName, workload.Cluster, workload.Namespace)
			}
		}
	}
	// delete application_resource
	if err = r.svcCtx.Database.WithContext(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "sqlite.DeleteApplicationResource")).Delete(&commontypes.ApplicationResource{}, "application_id = ?", id).Error; err != nil {
		return
	}
	r.Logger.Infof("Successfully delete application_resource for application=%s", application.AppName)
	// delete application_fass
	if err = r.svcCtx.Database.WithContext(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "sqlite.DeleteApplicationFaas")).Delete(&commontypes.ApplicationFaas{}, "application_id = ?", id).Error; err != nil {
		return
	}
	r.Logger.Infof("Successfully delete application_fass for application=%s", application.AppName)
	return
}

func (r *ApplicationRepository) Get(id, role string, tenant []string) (application *commontypes.Application, err error) {
	tx := r.svcCtx.Database.WithContext(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "sqlite.GetApplication")).Model(&commontypes.Application{})
	if role != constant.ROLE_ADMIN {
		tx.Where("tenant IN ?", tenant)
	}
	if err = tx.First(&application, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("error in GetApplication: %w", err)
	}
	// if err = tx.Preload("Template").First(&application, "id = ?", id).Error; err != nil {
	// 	return nil, fmt.Errorf("error in GetApplication: %w", err)
	// }
	return
}

func (r *ApplicationRepository) ListResource(req *commontypes.GetDataReq) (resp *types.Response, err error) {
	_, role, tenant, _ := utils.GetPayload(r.ctx)
	var data []commontypes.ApplicationResource
	tx := r.svcCtx.Database.WithContext(context.WithValue(r.ctx, commontypes.TraceIDKey{}, "sqlite.ListApplicationResource")).Model(&commontypes.ApplicationResource{})
	utils.SetTx(tx, req, nil, role, tenant, nil)
	if err = tx.Find(&data).Error; err != nil {
		r.Logger.Error(err)
		return
	}
	var wg sync.WaitGroup
	resultChan := make(chan *commontypes.ApplicationResource, len(data))
	for _, resource := range data {
		wg.Add(1)
		go func(resource commontypes.ApplicationResource, ch chan *commontypes.ApplicationResource) {
			resource.LiveManifest, err = r.taskService.GetResourceManifests(r.ctx, resource)
			if err != nil {
				r.Logger.Error(err)
				resource.LiveManifest = err.Error()
			}
			ch <- &resource
			wg.Done()
		}(resource, resultChan)
	}
	wg.Wait()
	close(resultChan)
	var res []commontypes.ApplicationResource
	for result := range resultChan {
		res = append(res, *result)
	}
	return &types.Response{
		Code: http.StatusOK,
		Data: res,
	}, nil
}

func (r *ApplicationRepository) saveIstio(appName, trafficPolicy string, workloads commontypes.WorkloadList, ifCreate bool) (err error) {
	cluster := workloads[0].Cluster
	namespace := workloads[0].Namespace
	var ag lizardagent.LizardAgent
	if ag, _, err = r.svcCtx.GetAgent(cluster, namespace); err != nil {
		return
	}
	if err = r.istioService.SaveDestinationRule(cluster, namespace, appName, workloads, ag, ifCreate); err != nil { // create destinationrule
		if strings.Contains(err.Error(), "not found") {
			r.istioService.SaveDestinationRule(cluster, namespace, appName, workloads, ag, true)
		}
	}
	if err = r.istioService.SaveVirtualService(cluster, namespace, appName, trafficPolicy, workloads, ag, ifCreate); err != nil { // create virtualservice
		if strings.Contains(err.Error(), "not found") {
			r.istioService.SaveVirtualService(cluster, namespace, appName, trafficPolicy, workloads, ag, true)
		}
	}
	return
}
