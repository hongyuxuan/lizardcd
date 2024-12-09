package svc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang-module/carbon"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/imroc/req/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

type DsService struct {
	logx.Logger
	ctx    context.Context
	svcCtx *ServiceContext
}

func NewDsService(ctx context.Context, svcCtx *ServiceContext) *DsService {
	return &DsService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (d *DsService) CreateProcessDefine(client *req.Client, name, projectCode string, params map[string]interface{}) (processCode int64, err error) {
	taskDefinitionJson, _ := json.Marshal(params["taskDefinitionJson"])
	locations, _ := json.Marshal(params["locations"])
	taskRelationJson, _ := json.Marshal(params["taskRelationJson"])
	var res commontypes.DolphinschedulerResponse
	if err = client.Post(fmt.Sprintf("/projects/%s/process-definition", projectCode)).
		SetFormData(map[string]string{
			"name":               name,
			"taskDefinitionJson": string(taskDefinitionJson),
			"locations":          string(locations),
			"taskRelationJson":   string(taskRelationJson),
		}).SetSuccessResult(&res).
		Do(context.WithValue(d.ctx, commontypes.TraceIDKey{}, "http.CreateDsProcessDefine")).Err; err != nil {
		return 0, fmt.Errorf("error in CreateDsProcessDefine: %w", err)
	}
	if res.Failed {
		return 0, fmt.Errorf("failed to CreateDsProcessDefine: %s", res.Msg)
	}
	data := res.Data.(map[string]interface{})
	processCode = int64(data["code"].(float64))
	d.Logger.Infof("Successfully CreateDsProcessDefine for %s, processCode=%d", name, processCode)
	return
}

func (d *DsService) ReleaseProcessDefine(client *req.Client, projectCode string, processCode int64, releaseState string) (err error) {
	var res commontypes.DolphinschedulerResponse
	if err = client.Post(fmt.Sprintf("/projects/%s/process-definition/%d/release", projectCode, processCode)).
		SetFormData(map[string]string{
			"releaseState": releaseState,
		}).SetSuccessResult(&res).
		Do(context.WithValue(d.ctx, commontypes.TraceIDKey{}, "http.ReleaseProcessDefine")).Err; err != nil {
		return fmt.Errorf("error in ReleaseProcessDefine: %w", err)
	}
	if res.Failed {
		return fmt.Errorf("failed to ReleaseProcessDefine %s: %s", releaseState, res.Msg)
	}
	d.Logger.Infof("Successfully ReleaseProcessDefine for processCode=%d", processCode)
	return
}

func (d *DsService) CreateScheduler(client *req.Client, projectCode string, processCode int64, cron string) (cronId int64, err error) {
	scheduler := map[string]string{
		"startTime":  carbon.Now().Format("Y-m-d H:i:s"),
		"endTime":    carbon.Now().AddYears(100).Format("Y-m-d H:i:s"),
		"crontab":    cron,
		"timezoneId": "Asia/Shanghai",
	}
	schedulerb, _ := json.Marshal(scheduler)
	var res commontypes.DolphinschedulerResponse
	if err = client.Post(fmt.Sprintf("/projects/%s/schedules", projectCode)).
		SetFormData(map[string]string{
			"processDefinitionCode": fmt.Sprintf("%d", processCode),
			"schedule":              string(schedulerb),
			"warningType":           "NONE",
			"failureStrategy":       "END",
		}).SetSuccessResult(&res).
		Do(context.WithValue(d.ctx, commontypes.TraceIDKey{}, "http.CreateScheduler")).Err; err != nil {
		return 0, fmt.Errorf("error in CreateScheduler: %w", err)
	}
	if res.Failed {
		return 0, fmt.Errorf("failed to CreateScheduler for processCode=%d: %s", processCode, res.Msg)
	}
	d.Logger.Infof("Successfully CreateScheduler for processCode=%d", processCode)
	return int64(res.Data.(map[string]interface{})["id"].(float64)), nil
}

func (d *DsService) ScheduleOnOffLine(client *req.Client, projectCode string, cronId int64, cronState string) (err error) {
	var res commontypes.DolphinschedulerResponse
	if err = client.Post(fmt.Sprintf("/projects/%s/schedules/%d/%s", projectCode, cronId, cronState)).
		Do(context.WithValue(d.ctx, commontypes.TraceIDKey{}, "http.ScheduleOnOffLine")).Err; err != nil {
		return fmt.Errorf("error in ScheduleOnOffLine: %w", err)
	}
	if res.Failed {
		return fmt.Errorf("failed to ScheduleOnOffLine for cronId=%d: %s", cronId, res.Msg)
	}
	d.Logger.Infof("Successfully ScheduleOnOffLine for cronId=%d", cronId)
	return
}

func (d *DsService) UpdateScheduler(client *req.Client, projectCode string, cronId int64, cron string) (err error) {
	var res commontypes.DolphinschedulerResponse
	scheduler := map[string]string{
		"startTime":  carbon.Now().Format("Y-m-d H:i:s"),
		"endTime":    carbon.Now().AddYears(100).Format("Y-m-d H:i:s"),
		"crontab":    cron,
		"timezoneId": "Asia/Shanghai",
	}
	schedulerb, _ := json.Marshal(scheduler)
	if err = client.Put(fmt.Sprintf("/projects/%s/schedules/%d", projectCode, cronId)).
		SetFormData(map[string]string{
			"schedule": string(schedulerb),
		}).SetSuccessResult(&res).
		Do(context.WithValue(d.ctx, commontypes.TraceIDKey{}, "http.UpdateScheduler")).Err; err != nil {
		return fmt.Errorf("error in UpdateScheduler: %w", err)
	}
	if res.Failed {
		return fmt.Errorf("failed to UpdateScheduler for cronId=%d: %s", cronId, res.Msg)
	}
	d.Logger.Infof("Successfully UpdateScheduler for cronId=%d", cronId)
	return
}

func (d *DsService) DeleteProcessDefine(client *req.Client, projectCode string, processCode int64) (err error) {
	var res commontypes.DolphinschedulerResponse
	// first offline
	if err = d.ReleaseProcessDefine(client, projectCode, processCode, "OFFLINE"); err != nil {
		return
	}
	if err = client.Delete(fmt.Sprintf("/projects/%s/process-definition/%d", projectCode, processCode)).
		SetSuccessResult(&res).
		Do(context.WithValue(d.ctx, commontypes.TraceIDKey{}, "http.DeleteProcessDefine")).Err; err != nil {
		return fmt.Errorf("error in DeleteProcessDefine: %w", err)
	}
	if res.Failed {
		return fmt.Errorf("failed to DeleteProcessDefine for processCode=%d: %s", processCode, res.Msg)
	}
	d.Logger.Infof("Successfully DeleteProcessDefine for processCode=%d", processCode)
	return
}
