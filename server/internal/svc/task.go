package svc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/golang-module/carbon"
	"github.com/google/uuid"
	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/constant"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	reqv3 "github.com/imroc/req/v3"
	"github.com/jinzhu/copier"
	"github.com/oliveagle/jsonpath"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel"
)

type TaskService struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *ServiceContext
	gitService *GitService
}

func NewTaskService(ctx context.Context, svcCtx *ServiceContext) *TaskService {
	return &TaskService{
		Logger:     logx.WithContext(ctx),
		ctx:        ctx,
		svcCtx:     svcCtx,
		gitService: NewGitService(ctx, svcCtx),
	}
}

type ResultChan struct {
	Cluster      string
	Namespace    string
	WorkloadType string
	WorkloadName string
	Success      bool
	Err          string
}

func (r ResultChan) ToString() string {
	return fmt.Sprintf("Cluster=%s Namespace=%s WorkloadType=%s WorkloadName=%s Success=%v Err=%s", r.Cluster, r.Namespace, r.WorkloadType, r.WorkloadName, r.Success, r.Err)
}

func (t *TaskService) RunTask(req *types.RunTaskReq, application commontypes.Application, tenant string) (taskId string, err error) {
	taskId = uuid.New().String()
	if req.Id != "" {
		taskId = req.Id
	}
	var repo commontypes.ImageRepository
	t.svcCtx.Sqlite.Model(&commontypes.ImageRepository{}).Where("id = ?", application.RepoId).First(&repo)
	// create task
	status := constant.TASK_STATUS_INITIALIZE
	if req.Waiting {
		status = constant.TASK_STATUS_WAITING
	}
	initAt := time.Now()
	if req.InitAt != "" {
		initAt = carbon.ParseByFormat(req.InitAt, "Y-m-d H:i:s").ToStdTime()
	}
	task := commontypes.TaskHistory{
		Id:          taskId,
		AppName:     req.AppName,
		TaskType:    req.TaskType,
		Status:      status,
		Tenant:      tenant,
		TriggerType: req.TriggerType,
		InitAt:      sql.NullTime{Time: initAt, Valid: true},
		Labels:      req.Labels,
	}
	if err = t.svcCtx.Sqlite.WithContext(context.WithValue(t.ctx, commontypes.TraceIDKey{}, "sqlite.SaveTaskHistory")).Save(&task).Error; err != nil {
		t.Logger.Error(err)
		return
	}
	res := t.svcCtx.Sqlite.WithContext(context.WithValue(t.ctx, commontypes.TraceIDKey{}, "sqlite.DeleteHistoryWorkload")).Where("task_history_id = ?", taskId).Delete(&commontypes.TaskHistoryWorkload{})
	t.Logger.Infof("Delete task_history_workload before task run, affect rows = %d", res.RowsAffected)

	// set workloads
	workloads := lo.Map(application.Workload, func(item commontypes.Workload, _ int) types.TaskWorkload {
		var wl types.TaskWorkload
		if item.Enable {
			copier.Copy(&wl, item)
			wl.ArtifactUrl = req.ArtifactUrl
		}
		return wl
	})
	if req.Workloads != nil {
		workloads = req.Workloads
	}

	if application.DeployType == "容器" {
		go t.ExecuteWorkload(task, workloads)
	} else if application.DeployType == "虚拟机" {
		go t.ExecuteVm(application, repo, task, workloads)
	} else if application.DeployType == "HTTP" {
		go t.ExecuteHttp(application, task, req.ArtifactUrl)
	} else if application.DeployType == "GitOps" {
		go t.ExecuteGitOps(application, task, workloads, tenant)
	}
	return
}

func (t *TaskService) ExecuteWorkload(task commontypes.TaskHistory, workloads []types.TaskWorkload) {
	var ag lizardagent.LizardAgent
	var err error
	var wg sync.WaitGroup
	wg.Add(len(workloads))
	results := make([]chan ResultChan, len(workloads))
	firstFail := false
	for i, w := range workloads {
		taskWorkload := commontypes.TaskHistoryWorkload{
			Workload: commontypes.Workload{
				Cluster:       w.Cluster,
				Namespace:     w.Namespace,
				WorkloadType:  w.WorkloadType,
				WorkloadName:  w.WorkloadName,
				ContainerName: w.ContainerName,
				ArtifactUrl:   w.ArtifactUrl,
			},
			TaskHistoryId: task.Id,
			UpdateAt:      time.Now(),
		}
		t.svcCtx.Sqlite.Create(&taskWorkload)
		// start to deploy
		if task.Status == constant.TASK_STATUS_WAITING {
			continue
		}
		if ag, err = t.svcCtx.GetAgent(w.Cluster, w.Namespace); err != nil {
			t.Logger.Error(err)
			return
		}
		if task.TaskType == constant.TASK_TYPE_DEPLOY {
			if w.WorkloadType == constant.K8S_RESOURCE_TYPE_DEPLOYMENTS {
				_, err = ag.PatchDeployment(context.Background(), &agent.PatchWorkloadRequest{
					Namespace:    w.Namespace,
					WorkloadName: w.WorkloadName,
					Container:    w.ContainerName,
					Image:        w.ArtifactUrl,
				})
			} else if w.WorkloadType == constant.K8S_RESOURCE_TYPE_STATEFULSETS {
				_, err = ag.PatchStatefulset(context.Background(), &agent.PatchWorkloadRequest{
					Namespace:    w.Namespace,
					WorkloadName: w.WorkloadName,
					Container:    w.ContainerName,
					Image:        w.ArtifactUrl,
				})
			}
		}
		if task.TaskType == constant.TASK_TYPE_ROLLOUT {
			if w.WorkloadType == constant.K8S_RESOURCE_TYPE_DEPLOYMENTS {
				_, err = ag.RolloutDeployment(context.Background(), &agent.RolloutWorkloadRequest{
					Namespace:    w.Namespace,
					WorkloadName: w.WorkloadName,
				})
			} else if w.WorkloadType == constant.K8S_RESOURCE_TYPE_STATEFULSETS {
				_, err = ag.RolloutStatefulset(context.Background(), &agent.RolloutWorkloadRequest{
					Namespace:    w.Namespace,
					WorkloadName: w.WorkloadName,
				})
			}
		}
		if err != nil {
			t.Logger.Error(err)
			firstFail = true
			// update task_history
			t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
				Status:     constant.TASK_STATUS_TERMINATED,
				Success:    sql.NullBool{Bool: false, Valid: true},
				ErrMessage: err.Error(),
				StartAt:    sql.NullTime{Time: time.Now(), Valid: true},
			})
			// update task_history_workload
			t.svcCtx.Sqlite.Model(&taskWorkload).Updates(commontypes.TaskHistoryWorkload{
				ErrMessage: err.Error(),
				UpdateAt:   time.Now(),
			})
			continue
		} else {
			t.Logger.Infof("Patch %s cluster=%s namespace=%s workload=%s container=%s image=%s", w.WorkloadType, w.Cluster, w.Namespace, w.WorkloadName, w.ContainerName, w.ArtifactUrl)
			if !firstFail {
				t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
					Status:  constant.TASK_STATUS_RUNNING,
					StartAt: sql.NullTime{Time: time.Now(), Valid: true},
				})
			}
		}
		// get workload status in background
		results[i] = make(chan ResultChan)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
		defer cancel()
		go t.getWorkloadStatus(ctx, taskWorkload, results[i], &wg)
	}
	if task.Status == constant.TASK_STATUS_WAITING {
		return
	}
	var failedWorkload []string
	for i := range workloads {
		var res ResultChan
		ch := results[i]
		go func(ch_comsume chan ResultChan) {
			res = <-ch_comsume
			if !res.Success {
				failedWorkload = append(failedWorkload, res.ToString())
			}
		}(ch)
	}
	wg.Wait()
	if len(failedWorkload) == 0 {
		t.Logger.Infof("Successfully run task, id=%s", task.Id)
		t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
			Status:   constant.TASK_STATUS_FINISHED,
			Success:  sql.NullBool{Bool: true, Valid: true},
			FinishAt: sql.NullTime{Time: time.Now(), Valid: true},
			Expire:   computeExpire(task.StartAt),
		})
	} else {
		t.Logger.Errorf("Failed run task, id=%s", task.Id)
		failB, _ := json.Marshal(failedWorkload)
		t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
			Status:     constant.TASK_STATUS_FINISHED,
			Success:    sql.NullBool{Bool: false, Valid: true},
			ErrMessage: string(failB),
			FinishAt:   sql.NullTime{Time: time.Now(), Valid: true},
			Expire:     computeExpire(task.StartAt),
		})
	}
}

func (t *TaskService) ExecuteVm(application commontypes.Application, repo commontypes.ImageRepository, task commontypes.TaskHistory, workloads []types.TaskWorkload) {
	var req types.VmDeployReq
	json.Unmarshal([]byte(application.ExtraVars), &req)
	var err error
	var wg sync.WaitGroup
	wg.Add(len(workloads))
	results := make([]chan ResultChan, len(workloads))
	firstFail := false
	for i, w := range workloads {
		taskWorkload := commontypes.TaskHistoryWorkload{
			Workload: commontypes.Workload{
				WorkloadType: w.WorkloadType,
				WorkloadName: w.WorkloadName,
				ArtifactUrl:  w.ArtifactUrl,
			},
			TaskHistoryId: task.Id,
			UpdateAt:      time.Now(),
		}
		t.svcCtx.Sqlite.Create(&taskWorkload)
		// start to deploy
		if task.Status == constant.TASK_STATUS_WAITING {
			continue
		}
		req.ArtifactUrl = w.ArtifactUrl
		req.ArtifactHeader = map[string]string{"X-JFrog-Art-Api": repo.RepoPassword}
		req.Targets = []string{w.WorkloadName}
		_, err = t.VmDeploy(context.Background(), &req)
		if err != nil {
			t.Logger.Error(err)
			firstFail = true
			// update task_history
			t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
				Success:    sql.NullBool{Bool: false, Valid: true},
				ErrMessage: err.Error(),
				StartAt:    sql.NullTime{Time: time.Now(), Valid: true},
			})
			// update task_history_workload
			t.svcCtx.Sqlite.Model(&taskWorkload).Updates(commontypes.TaskHistoryWorkload{
				ErrMessage: err.Error(),
				UpdateAt:   time.Now(),
			})
			continue
		} else {
			t.Logger.Infof("Execute vm deploy on host=%s", w.WorkloadName)
			if !firstFail {
				t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
					Status:  constant.TASK_STATUS_RUNNING,
					StartAt: sql.NullTime{Time: time.Now(), Valid: true},
				})
			}
		}
		// get workload status in background
		results[i] = make(chan ResultChan)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
		defer cancel()
		go t.getVmStatus(ctx, taskWorkload, &req.HealthCheck, results[i], &wg)
	}
	if task.Status == constant.TASK_STATUS_WAITING {
		return
	}
	var failedWorkload []string
	for i := range workloads {
		var res ResultChan
		ch := results[i]
		go func(ch_comsume chan ResultChan) {
			res = <-ch_comsume
			if !res.Success {
				failedWorkload = append(failedWorkload, res.ToString())
			}
		}(ch)
	}
	wg.Wait()
	if len(failedWorkload) == 0 {
		t.Logger.Infof("Successfully run task, id=%s", task.Id)
		t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
			Status:   constant.TASK_STATUS_FINISHED,
			Success:  sql.NullBool{Bool: true, Valid: true},
			FinishAt: sql.NullTime{Time: time.Now(), Valid: true},
			Expire:   computeExpire(task.StartAt),
		})
	} else {
		t.Logger.Errorf("Failed run task, id=%s", task.Id)
		failB, _ := json.Marshal(failedWorkload)
		t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
			Status:     constant.TASK_STATUS_FINISHED,
			Success:    sql.NullBool{Bool: false, Valid: true},
			ErrMessage: string(failB),
			FinishAt:   sql.NullTime{Time: time.Now(), Valid: true},
			Expire:     computeExpire(task.StartAt),
		})
	}
}

func (t *TaskService) ExecuteHttp(application commontypes.Application, task commontypes.TaskHistory, artifactUrl string) {
	var req types.HttpDeployReq
	json.Unmarshal([]byte(application.ExtraVars), &req)
	req.ArtifactUrl = artifactUrl
	// update task_history_workload
	taskWorkload := commontypes.TaskHistoryWorkload{
		Workload: commontypes.Workload{
			WorkloadType: application.DeployType,
			WorkloadName: req.HttpUrl,
			ArtifactUrl:  artifactUrl,
		},
		TaskHistoryId: task.Id,
		UpdateAt:      time.Now(),
	}
	t.svcCtx.Sqlite.Create(&taskWorkload)
	if task.Status == constant.TASK_STATUS_WAITING {
		return
	}

	// var res *types.Response
	res, err := t.HttpDeploy(context.Background(), &req)
	if err != nil {
		// update task_history
		t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
			Status:     constant.TASK_STATUS_TERMINATED,
			Success:    sql.NullBool{Bool: false, Valid: true},
			ErrMessage: err.Error(),
			StartAt:    sql.NullTime{Time: time.Now(), Valid: true},
		})
		// update task_history_workload
		t.svcCtx.Sqlite.Model(&taskWorkload).Updates(commontypes.TaskHistoryWorkload{
			ErrMessage: err.Error(),
			UpdateAt:   time.Now(),
		})
		return
	}
	re, _ := regexp.Compile(`.*(\{\{response(\$.*)\}\}).*`)
	matches := re.FindStringSubmatch(req.HealthCheck.HttpPath)
	if len(matches) >= 3 {
		var lookup interface{}
		if lookup, err = jsonpath.JsonPathLookup(res, matches[2]); err != nil {
			t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
				Status:     constant.TASK_STATUS_TERMINATED,
				ErrMessage: err.Error(),
				StartAt:    sql.NullTime{Time: time.Now(), Valid: true},
			})
		}
		lookupstr := utils.AnyToString(lookup)
		req.HealthCheck.HttpPath = strings.ReplaceAll(req.HealthCheck.HttpPath, matches[1], lookupstr)
		req.HealthCheck.HttpBody = strings.ReplaceAll(req.HealthCheck.HttpBody, matches[1], lookupstr)
	}
	// update task_workload
	t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
		Status:  constant.TASK_STATUS_RUNNING,
		StartAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	t.Logger.Infof("Successfully start http task, response: %s", res)

	// httpcheck in background
	ctx, _ := context.WithTimeout(context.Background(), time.Second*300)
	// defer cancel()
	go t.getHttpStatus(ctx, taskWorkload, &req.HealthCheck, req.HttpUrl, req.HttpHeader, task)
}

func (t *TaskService) ExecuteGitOps(application commontypes.Application, task commontypes.TaskHistory, workloads []types.TaskWorkload, tenant string) {
	var ag lizardagent.LizardAgent
	var err error
	var wg sync.WaitGroup
	wg.Add(len(workloads))
	firstFail := false
	// get yaml from gitlab
	var commitId, yamlstring string
	if commitId, yamlstring, err = t.gitService.GetYamlFromGit(application, tenant); err != nil {
		t.Logger.Error(err)
		t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
			Status:     constant.TASK_STATUS_TERMINATED,
			Success:    sql.NullBool{Bool: false, Valid: true},
			ErrMessage: err.Error(),
			StartAt:    sql.NullTime{Time: time.Now(), Valid: true},
		})
	}
	// parse yaml
	if err = t.gitService.ParseYaml(context.Background(), application, yamlstring); err != nil {
		t.Logger.Error(err)
		t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
			Status:     constant.TASK_STATUS_TERMINATED,
			Success:    sql.NullBool{Bool: false, Valid: true},
			ErrMessage: err.Error(),
			StartAt:    sql.NullTime{Time: time.Now(), Valid: true},
		})
	}
	// send to agent
	var appResource []commontypes.ApplicationResource
	t.svcCtx.Sqlite.Model(&commontypes.ApplicationResource{}).Where("application_id = ?", application.Id).Find(&appResource)
	var taskWorkloads []commontypes.TaskHistoryWorkload
	for _, w := range workloads {
		for _, ar := range appResource {
			if ar.ResourceType == constant.K8S_RESOURCE_TYPE_DEPLOYMENTS || ar.ResourceType == constant.K8S_RESOURCE_TYPE_STATEFULSETS {
				taskWorkload := commontypes.TaskHistoryWorkload{
					Workload: commontypes.Workload{
						Cluster:      w.Cluster,
						Namespace:    w.Namespace,
						WorkloadType: ar.ResourceType,
						WorkloadName: ar.ResourceName,
					},
					TaskHistoryId: task.Id,
					UpdateAt:      time.Now(),
				}
				t.svcCtx.Sqlite.Create(&taskWorkload)
				taskWorkloads = append(taskWorkloads, taskWorkload)
			}
		}
	}
	for _, w := range workloads {
		if ag, err = t.svcCtx.GetAgent(w.Cluster, w.Namespace); err != nil {
			t.Logger.Error(err)
			return
		}
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.ApplyYaml(context.Background(), &agent.YamlRequest{
			Namespace: w.Namespace,
			Ymlstring: yamlstring,
		}); err != nil {
			t.Logger.Error(err)
			firstFail = true
			// update task_history
			t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
				Status:     constant.TASK_STATUS_TERMINATED,
				Success:    sql.NullBool{Bool: false, Valid: true},
				ErrMessage: err.Error(),
				StartAt:    sql.NullTime{Time: time.Now(), Valid: true},
			})
			continue
		} else {
			var res struct {
				WorkloadType string
				WorkloadName string
			}
			json.Unmarshal(rpcResponse.Data, &res)
			t.Logger.Infof("Synchronize send yaml to cluster=%s namespace=%s success", w.Cluster, w.Namespace)
			if !firstFail {
				t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
					Status:  constant.TASK_STATUS_RUNNING,
					StartAt: sql.NullTime{Time: time.Now(), Valid: true},
				})
			}
		}
	}
	results := make([]chan ResultChan, len(taskWorkloads))
	for i, thw := range taskWorkloads {
		// get workload status in background
		results[i] = make(chan ResultChan)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
		defer cancel()
		go t.getWorkloadStatus(ctx, thw, results[i], &wg)
	}
	// update application
	application.GitOps.CommitId = commitId
	t.svcCtx.Sqlite.Model(&commontypes.Application{}).Where("id = ?", application.Id).Updates(&commontypes.Application{
		GitOps: application.GitOps,
	})
	var failedWorkload []string
	for i := range workloads {
		var res ResultChan
		ch := results[i]
		go func(ch_comsume chan ResultChan) {
			res = <-ch_comsume
			if !res.Success {
				failedWorkload = append(failedWorkload, res.ToString())
			}
		}(ch)
	}
	wg.Wait()
	if len(failedWorkload) == 0 {
		t.Logger.Infof("Successfully run task, id=%s", task.Id)
		t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
			Status:   constant.TASK_STATUS_FINISHED,
			Success:  sql.NullBool{Bool: true, Valid: true},
			FinishAt: sql.NullTime{Time: time.Now(), Valid: true},
			Expire:   computeExpire(task.StartAt),
		})
		// get last manifests of application_resource
		for _, resource := range appResource {
			manifests, err := t.GetResourceManifests(context.Background(), resource)
			if err != nil {
				t.Logger.Error(err)
			}
			if err = t.svcCtx.Sqlite.Model(&commontypes.ApplicationResource{}).
				Where("application_id = ?", resource.ApplicationId).
				Where("cluster = ?", resource.Cluster).
				Where("namespace = ?", resource.Namespace).
				Where("resource_type = ?", resource.ResourceType).
				Where("resource_name = ?", resource.ResourceName).Update("last_manifest", manifests).Error; err != nil {
				t.Logger.Error(err)
			}
		}
	} else {
		t.Logger.Errorf("Failed run task, id=%s", task.Id)
		failB, _ := json.Marshal(failedWorkload)
		t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
			Status:     constant.TASK_STATUS_FINISHED,
			Success:    sql.NullBool{Bool: false, Valid: true},
			ErrMessage: string(failB),
			FinishAt:   sql.NullTime{Time: time.Now(), Valid: true},
			Expire:     computeExpire(task.StartAt),
		})
	}
}

func (t *TaskService) VmDeploy(ctx context.Context, req *types.VmDeployReq) (result map[string]string, err error) {
	var ag lizardagent.LizardAgent
	result = make(map[string]string)
	for _, ip := range req.Targets {
		if ag, err = t.svcCtx.GetTargetAgent(ip); err != nil {
			return
		}
		var rpcResponse *agent.Response
		header, _ := json.Marshal(req.ArtifactHeader)
		if rpcResponse, err = ag.VmDeploy(context.WithValue(ctx, commontypes.TraceIDKey{}, "rpc.VmDeploy"), &agent.VmDeployRequest{
			ArtifactUrl:    req.ArtifactUrl,
			ArtifactHeader: header,
			DeployPath:     req.DeployPath,
			DeployUser:     req.DeployUser,
			PreCommand:     req.PreCommand,
			StartCommand:   req.StartCommand,
		}); err != nil {
			t.Logger.Error(err)
			return
		}
		result[ip] = string(rpcResponse.Data)
	}
	return
}

func (t *TaskService) Healthcheck(ctx context.Context, req *types.HealthCheck, target string) (message string, data interface{}, err error) {
	var ag lizardagent.LizardAgent
	if ag, err = t.svcCtx.GetTargetAgent(target); err != nil {
		return
	}
	if req.Type == "none" {
		message = target + " do not have healthcheck"
	} else {
		var rpcResponse *agent.Response
		if rpcResponse, err = ag.VmHealthCheck(context.WithValue(ctx, commontypes.TraceIDKey{}, "rpc.VmHealthCheck"), &agent.VmHealthCheckRequest{
			Type:   req.Type,
			Method: req.Method,
			Port:   req.Port,
			Uri:    req.Uri,
			Shell:  req.Shell,
		}); err != nil {
			t.Logger.Error(err)
			return
		}
		data = string(rpcResponse.Data)
		message = target + " healthcheck success"
	}
	return
}

func (t *TaskService) HttpDeploy(ctx context.Context, req *types.HttpDeployReq) (resInterface interface{}, err error) {
	req.HttpBody = strings.ReplaceAll(req.HttpBody, "{{artifact_url}}", req.ArtifactUrl)
	client := utils.NewHttpClient(otel.Tracer("imroc/req"))
	if t.svcCtx.Config.Log.Level == "debug" {
		client.EnableDebug(true)
	}
	request := client.SetBaseURL(req.HttpUrl).SetCommonHeaders(req.HttpHeader).R()
	if req.HttpContentType == "json" {
		var httpBody map[string]interface{}
		if err = json.Unmarshal([]byte(req.HttpBody), &httpBody); err != nil {
			t.Logger.Error(err)
			return
		}
		request.SetBody(httpBody)
	} else if req.HttpContentType == "x-www-form-urlencoded" {
		var httpBody map[string]string
		if err = json.Unmarshal([]byte(req.HttpBody), &httpBody); err != nil {
			t.Logger.Error(err)
			return
		}
		request.SetFormData(httpBody)
	} else {
		return nil, errorx.NewDefaultError("Unsupported content-type: application/%s", req.HttpContentType)
	}
	var res *reqv3.Response
	if req.HttpMethod == "post" {
		res, err = request.Post(req.HttpPath)
	} else if req.HttpMethod == "put" {
		res, err = request.Put(req.HttpPath)
	} else {
		return nil, errorx.NewDefaultError("Unsupported http method: %s", req.HttpMethod)
	}
	if err != nil {
		t.Logger.Error(err)
		return
	}
	if res.IsError() {
		e := errorx.NewDefaultError("Http deploy return statusCode: %d, data: %v", res.StatusCode, res.String())
		t.Logger.Error(e)
		return nil, e
	}
	res.Unmarshal(&resInterface)
	if req.ResJsonpath != "" {
		var lookup interface{}
		if lookup, err = jsonpath.JsonPathLookup(resInterface, req.ResJsonpath); err != nil {
			t.Logger.Error(err)
			return
		}
		lookupstr := utils.AnyToString(lookup)
		re, _ := regexp.Compile(req.ResKeyword)
		if !re.MatchString(lookupstr) {
			return nil, errorx.NewDefaultError("Http deploy failed with response: %s", res.String())
		}
	}
	return
}

func (t *TaskService) HttpCheck(ctx context.Context, req *types.HttpCheck, httpUrl string, httpHeader map[string]string) (finished, success bool, message string, err error) {
	client := utils.NewHttpClient(otel.Tracer("imroc/req"))
	if t.svcCtx.Config.Log.Level == "debug" {
		client.EnableDebug(true)
	}
	request := client.SetBaseURL(httpUrl).SetCommonHeaders(httpHeader).R()
	var res *reqv3.Response
	if req.Method == "get" {
		res, err = request.Get(req.HttpPath)
	} else if req.Method == "post" {
		var httpBody map[string]string
		if err = json.Unmarshal([]byte(req.HttpBody), &httpBody); err != nil {
			t.Logger.Error(err)
			return
		}
		res, err = request.SetBody(httpBody).Post(req.HttpPath)
	} else { // no health check
		return
	}
	if err != nil {
		t.Logger.Error(err)
		return
	}
	if res.IsError() {
		err = errorx.NewDefaultError("Http check return statusCode: %d, data: %v", res.StatusCode, res.String())
		t.Logger.Error(err)
		return
	}

	var resInterface interface{}
	res.Unmarshal(&resInterface)
	t.Logger.Debugf("Http check response: %s", res.String())
	finished = true
	success = true
	var lookup interface{}
	if req.FinishJsonpath != "" {
		if lookup, err = jsonpath.JsonPathLookup(resInterface, req.FinishJsonpath); err != nil {
			t.Logger.Error(err)
			return
		}
		lookupstr := utils.AnyToString(lookup)
		re, _ := regexp.Compile(req.FinishKeyword)
		if !re.MatchString(lookupstr) {
			finished = false
		}
	}
	if req.SuccessJsonpath != "" {
		if lookup, err = jsonpath.JsonPathLookup(resInterface, req.SuccessJsonpath); err != nil {
			t.Logger.Error(err)
			return
		}
		lookupstr := utils.AnyToString(lookup)
		re, _ := regexp.Compile(req.SuccessKeyword)
		if !re.MatchString(lookupstr) {
			success = false
		}
	}
	if req.MsgJsonpath != "" {
		re, _ := regexp.Compile(`.*(\{\{(\$.*)\}\}).*`)
		matches := re.FindStringSubmatch(req.MsgJsonpath)
		if len(matches) >= 3 {
			var lookup interface{}
			if lookup, err = jsonpath.JsonPathLookup(resInterface, matches[2]); err != nil {
				t.Logger.Error(err)
				return
			}
			lookupstr := utils.AnyToString(lookup)
			message = strings.ReplaceAll(req.MsgJsonpath, matches[1], lookupstr)
		}
	}
	return
}

func (t *TaskService) GetResourceManifests(ctx context.Context, resource commontypes.ApplicationResource) (manifests string, err error) {
	var rpcResponse *agent.YamlResponse
	var ag lizardagent.LizardAgent
	if ag, err = t.svcCtx.GetAgent(resource.Cluster, resource.Namespace); err != nil {
		return
	}
	if rpcResponse, err = ag.Getyaml(ctx, &lizardagent.GetYamlRequest{
		Namespace:    resource.Namespace,
		ResourceType: resource.ResourceType,
		ResourceName: resource.ResourceName,
	}); err != nil {
		return "", fmt.Errorf("error in Getyaml from agent: %w", err)
	}
	return rpcResponse.Data, nil
}

func (t *TaskService) getWorkloadStatus(ctx context.Context, taskWorkload commontypes.TaskHistoryWorkload, result chan ResultChan, wg *sync.WaitGroup) {
	var ag lizardagent.LizardAgent
	var err error
	// sleep 10s, waiting for kubernetes
	time.Sleep(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			t.Logger.Errorf("Cluster=%s namespace=%s workload_type=%s workload_name=%s running TIMEOUT and TERMINATED", taskWorkload.Workload.Cluster, taskWorkload.Workload.Namespace, taskWorkload.Workload.WorkloadType, taskWorkload.Workload.WorkloadName)
			t.setStatus(taskWorkload, false, nil, errorx.NewDefaultError("TIMEOUT and TERMINATED"), result)
			time.Sleep(1 * time.Second)
			wg.Done()
			return
		default:
			podStatus := true
			if ag, err = t.svcCtx.GetAgent(taskWorkload.Workload.Cluster, taskWorkload.Workload.Namespace); err != nil {
				t.Logger.Error(err)
				t.setStatus(taskWorkload, false, nil, err, result)
				return
			}
			var rpcResponse *agent.Response
			if rpcResponse, err = ag.GetPodStatus(context.Background(), &agent.GetWorkloadRequest{
				Namespace:    taskWorkload.Workload.Namespace,
				WorkloadName: taskWorkload.Workload.WorkloadName,
				WorkloadType: taskWorkload.Workload.WorkloadType,
			}); err != nil {
				t.Logger.Error(err)
				t.setStatus(taskWorkload, false, nil, err, result)
				return
			}
			var status *commontypes.WorkloadStatus
			json.Unmarshal(rpcResponse.Data, &status)
			t.Logger.Infof("Cluster=%s namespace=%s workload_type=%s workload_name=%s pod status=%v", taskWorkload.Workload.Cluster, taskWorkload.Workload.Namespace, taskWorkload.Workload.WorkloadType, taskWorkload.Workload.WorkloadName, status.Pods)
			for _, pod := range status.Pods {
				if pod.Ready == "False" {
					podStatus = false
				}
			}
			// status write to database
			b, _ := json.Marshal(status.Pods)
			t.svcCtx.Sqlite.Model(&taskWorkload).Updates(commontypes.TaskHistoryWorkload{
				Status:   string(b),
				UpdateAt: time.Now(),
			})
			if podStatus { // 任务结束
				t.setStatus(taskWorkload, true, status, nil, result)
				wg.Done()
				break
			}
			time.Sleep(3 * time.Second)
		}
	}
}

func (t *TaskService) getVmStatus(ctx context.Context, taskWorkload commontypes.TaskHistoryWorkload, healthCheck *types.HealthCheck, result chan ResultChan, wg *sync.WaitGroup) {
	time.Sleep(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			t.Logger.Errorf("Vm host=%s running TIMEOUT and TERMINATED", taskWorkload.Workload.WorkloadName)
			t.setStatus(taskWorkload, false, nil, errorx.NewDefaultError("TIMEOUT and TERMINATED"), result)
			time.Sleep(1 * time.Second)
			wg.Done()
			return
		default:
			if message, data, err := t.Healthcheck(context.Background(), &types.HealthCheck{
				Type:   healthCheck.Type,
				Method: healthCheck.Method,
				Port:   healthCheck.Port,
				Uri:    healthCheck.Uri,
				Shell:  healthCheck.Shell,
			}, taskWorkload.Workload.WorkloadName); err != nil {
				t.Logger.Infof("Vm host healthcheck failed: %v, output: %v", err, data)
				t.svcCtx.Sqlite.Model(&taskWorkload).Updates(commontypes.TaskHistoryWorkload{
					Status:   taskWorkload.Workload.WorkloadName + ": " + err.Error(),
					UpdateAt: time.Now(),
				})
			} else {
				t.Logger.Info(message)
				t.setStatus(taskWorkload, true, nil, nil, result)
				time.Sleep(1 * time.Second)
				wg.Done()
				break
			}
			time.Sleep(3 * time.Second)
		}
	}
}

func (t *TaskService) getHttpStatus(ctx context.Context, taskWorkload commontypes.TaskHistoryWorkload, httpCheck *types.HttpCheck, httpUrl string, httpHeader map[string]string, task commontypes.TaskHistory) {
	time.Sleep(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			e := errorx.NewDefaultError("Http task running TIMEOUT and TERMINATED")
			t.Logger.Error(e)
			t.setHttpStatus(task, taskWorkload, false, e.Error())
			return
		default:
			finished, success, message, err := t.HttpCheck(context.Background(), httpCheck, httpUrl, httpHeader)
			if err != nil {
				t.Logger.Infof("Failed run http task: %v", err)
				t.setHttpStatus(task, taskWorkload, false, err.Error())
				return
			} else {
				if !finished {
					t.Logger.Infof("Http task is running, finished=%v success=%v message=%v", finished, success, message)
				} else {
					t.Logger.Infof("Http task finished, finished=%v success=%v message=%v", finished, success, message)
					t.setHttpStatus(task, taskWorkload, success, message)
					return
				}
			}
			time.Sleep(3 * time.Second)
		}
	}
}

func (t *TaskService) setStatus(taskWorkload commontypes.TaskHistoryWorkload, success bool, status *commontypes.WorkloadStatus, err error, ch chan ResultChan) {
	thw := commontypes.TaskHistoryWorkload{
		UpdateAt: time.Now(),
	}
	r := ResultChan{
		Cluster:      taskWorkload.Workload.Cluster,
		Namespace:    taskWorkload.Workload.Namespace,
		WorkloadType: taskWorkload.Workload.WorkloadType,
		WorkloadName: taskWorkload.Workload.WorkloadName,
		Success:      success,
	}
	if status != nil {
		b, _ := json.Marshal(status.Pods)
		// taskWorkload.Workload.Revision = status.Revision
		thw.Workload = taskWorkload.Workload
		thw.Status = string(b)
	}
	if err != nil {
		thw.ErrMessage = err.Error()
		r.Err = err.Error()
	}
	t.svcCtx.Sqlite.Model(&taskWorkload).Updates(thw)
	ch <- r
}

func (t *TaskService) setHttpStatus(task commontypes.TaskHistory, taskWorkload commontypes.TaskHistoryWorkload, success bool, errMessage string) {
	// update task_workload
	t.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
		Status:     constant.TASK_STATUS_FINISHED,
		Success:    sql.NullBool{Bool: success, Valid: true},
		ErrMessage: errMessage,
		FinishAt:   sql.NullTime{Time: time.Now(), Valid: true},
		Expire:     computeExpire(task.StartAt),
	})
	// update task_history_workload
	t.svcCtx.Sqlite.Model(&taskWorkload).Updates(commontypes.TaskHistoryWorkload{
		ErrMessage: errMessage,
		UpdateAt:   time.Now(),
	})
}

func computeExpire(t sql.NullTime) (expire string) {
	if t.Valid {
		expire = time.Since(t.Time).Truncate(time.Duration(1) * time.Millisecond).String()
	}
	return
}
