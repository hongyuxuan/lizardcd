package task

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/lizardagent"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/hongyuxuan/lizardcd/common/constant"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"github.com/hongyuxuan/lizardcd/server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/google/uuid"
)

type RunTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
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

func NewRunTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RunTaskLogic {
	return &RunTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RunTaskLogic) RunTask(req *types.RunTaskRepoReq) (resp *types.Response, err error) {
	_, _, tenant, _ := utils.GetPayload(l.ctx)
	id := uuid.New().String()
	if req.Id != "" {
		id = req.Id
	}
	task := commontypes.TaskHistory{
		Id:          id,
		AppName:     req.AppName,
		TaskType:    req.TaskType,
		Status:      "initialize",
		Tenant:      tenant,
		TriggerType: req.TriggerType,
		InitAt:      sql.NullTime{Time: time.Now(), Valid: true},
		Labels:      req.Labels,
	}
	if err = l.svcCtx.Sqlite.Save(&task).Error; err != nil {
		l.Logger.Error(err)
		return
	}

	go l.getStatus(task, req.Workloads)

	resp = &types.Response{
		Code:    http.StatusOK,
		Message: "任务提交成功",
		Data: map[string]string{
			"id": task.Id,
		},
	}
	return
}

func (l *RunTaskLogic) getStatus(task commontypes.TaskHistory, workloads []types.TaskWorkload) {
	var ag lizardagent.LizardAgent
	var err error
	var wg sync.WaitGroup
	wg.Add(len(workloads))
	results := make([]chan ResultChan, len(workloads))
	firstFail := false
	for i, w := range workloads {
		taskWorkload := commontypes.TaskHistoryWorkload{
			Workload: commontypes.WorkLoad{
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
		l.svcCtx.Sqlite.Create(&taskWorkload)
		// start to deploy
		if ag, err = l.svcCtx.GetAgent(w.Cluster, w.Namespace); err != nil {
			l.Logger.Error(err)
			return
		}
		if task.TaskType == constant.K8S_TASK_TYPE_DEPLOY {
			_, err = ag.PatchDeployment(context.Background(), &agent.PatchWorkloadRequest{
				Namespace:    w.Namespace,
				WorkloadName: w.WorkloadName,
				Container:    w.ContainerName,
				Image:        w.ArtifactUrl,
			})
		}
		if task.TaskType == constant.K8S_TASK_TYPE_ROLLOUT {
			_, err = ag.RolloutDeployment(context.Background(), &agent.RolloutWorkloadRequest{
				Namespace:    w.Namespace,
				WorkloadName: w.WorkloadName,
			})
		}
		if err != nil {
			l.Logger.Error(err)
			firstFail = true
			// update task_history
			l.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
				Success:    sql.NullBool{Bool: false, Valid: true},
				ErrMessage: err.Error(),
				StartAt:    sql.NullTime{Time: time.Now(), Valid: true},
			})
			// update task_history_workload
			l.svcCtx.Sqlite.Model(&taskWorkload).Updates(commontypes.TaskHistoryWorkload{
				ErrMessage: err.Error(),
				UpdateAt:   time.Now(),
			})
			continue
		} else {
			l.Logger.Infof("Patch deployments cluster=%s namespace=%s workload=%s container=%s image=%s", w.Cluster, w.Namespace, w.WorkloadName, w.ContainerName, w.ArtifactUrl)
			if !firstFail {
				l.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
					Status:  "running",
					StartAt: sql.NullTime{Time: time.Now(), Valid: true},
				})
			}
		}
		// get workload status in background
		results[i] = make(chan ResultChan)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
		defer cancel()
		go l.getWorkloadStatus(ctx, taskWorkload, results[i], &wg)
	}
	var failedWorkload []string
	for i := range workloads {
		var res ResultChan
		ch := results[i]
		go func(ch_comsume chan ResultChan) {
			res = <-ch_comsume
			if res.Success == false {
				failedWorkload = append(failedWorkload, res.ToString())
			}
		}(ch)
	}
	wg.Wait()
	if len(failedWorkload) == 0 {
		l.Logger.Infof("Successfully run task, id=%s", task.Id)
		l.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
			Status:   "finished",
			Success:  sql.NullBool{Bool: true, Valid: true},
			FinishAt: sql.NullTime{Time: time.Now(), Valid: true},
			Expire:   time.Now().Sub(task.StartAt.Time).Truncate(time.Duration(1) * time.Millisecond).String(),
		})
	} else {
		l.Logger.Errorf("Failed run task, id=%s", task.Id)
		failB, _ := json.Marshal(failedWorkload)
		l.svcCtx.Sqlite.Model(&task).Updates(commontypes.TaskHistory{
			Status:     "finished",
			Success:    sql.NullBool{Bool: false, Valid: true},
			ErrMessage: string(failB),
			FinishAt:   sql.NullTime{Time: time.Now(), Valid: true},
			Expire:     time.Now().Sub(task.StartAt.Time).Truncate(time.Duration(1) * time.Millisecond).String(),
		})
	}
}

func (l *RunTaskLogic) getWorkloadStatus(ctx context.Context, taskWorkload commontypes.TaskHistoryWorkload, result chan ResultChan, wg *sync.WaitGroup) {
	var ag lizardagent.LizardAgent
	var err error
	// sleep 10s, waiting for kubernetes
	time.Sleep(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			l.Logger.Errorf("Cluster=%s namespace=%s workload_type=%s workload_name=%s running TIMEOUT and TERMINATED", taskWorkload.Workload.Cluster, taskWorkload.Workload.Namespace, taskWorkload.Workload.WorkloadType, taskWorkload.Workload.WorkloadName)
			l.setStatus(taskWorkload, false, nil, fmt.Errorf("TIMEOUT and TERMINATED"), result)
			wg.Done()
			break
		default:
			podStatus := true
			if ag, err = l.svcCtx.GetAgent(taskWorkload.Workload.Cluster, taskWorkload.Workload.Namespace); err != nil {
				l.Logger.Error(err)
				return
			}
			var rpcResponse *agent.Response
			if rpcResponse, err = ag.GetPodStatus(context.Background(), &agent.GetWorkloadRequest{
				Namespace:    taskWorkload.Workload.Namespace,
				WorkloadName: taskWorkload.Workload.WorkloadName,
			}); err != nil {
				l.setStatus(taskWorkload, false, nil, err, result)
				continue
			}
			var status *commontypes.WorkloadStatus
			json.Unmarshal(rpcResponse.Data, &status)
			l.Logger.Infof("Cluster=%s namespace=%s workload_type=%s workload_name=%s pod status=%v", taskWorkload.Workload.Cluster, taskWorkload.Workload.Namespace, taskWorkload.Workload.WorkloadType, taskWorkload.Workload.WorkloadName, status.Pods)
			for _, pod := range status.Pods {
				if pod.Ready == "False" {
					podStatus = false
				}
			}
			// status write to database
			b, _ := json.Marshal(status.Pods)
			l.svcCtx.Sqlite.Model(&taskWorkload).Updates(commontypes.TaskHistoryWorkload{
				Status:   string(b),
				UpdateAt: time.Now(),
			})
			if podStatus == true { // 任务结束
				l.setStatus(taskWorkload, true, status.Pods, nil, result)
				wg.Done()
				break
			}
			time.Sleep(3 * time.Second)
		}
	}
}

func (l *RunTaskLogic) setStatus(taskWorkload commontypes.TaskHistoryWorkload, success bool, pods []commontypes.PodStatus, err error, ch chan ResultChan) {
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
	if pods != nil {
		status, _ := json.Marshal(pods)
		thw.Status = string(status)
	}
	if err != nil {
		thw.ErrMessage = err.Error()
		r.Err = err.Error()
	}
	l.svcCtx.Sqlite.Model(&taskWorkload).Updates(thw)
	ch <- r
}
