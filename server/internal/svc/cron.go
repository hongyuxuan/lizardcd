package svc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"html/template"

	"github.com/hongyuxuan/lizardcd/common/constant"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/hongyuxuan/lizardcd/server/internal/types"
	"github.com/imroc/req/v3"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel"
)

type CronService struct {
	logx.Logger
	ctx         context.Context
	svcCtx      *ServiceContext
	taskService *TaskService
	dsService   *DsService
}

func NewCronService(ctx context.Context, svcCtx *ServiceContext) *CronService {
	return &CronService{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		taskService: NewTaskService(context.Background(), svcCtx),
		dsService:   NewDsService(context.Background(), svcCtx),
	}
}

func (c *CronService) RemoveCron(application commontypes.Application, tenant string, ifDelete bool) {
	if _, ok := c.svcCtx.CronIdMap[application.AppName]; ok {
		c.svcCtx.Cron.Remove(c.svcCtx.CronIdMap[application.AppName].CronId)
		delete(c.svcCtx.CronIdMap, application.AppName)
	}
	if application.GitOps != nil && (application.GitOps.SyncType == "手动" || !application.GitOps.UseDS || ifDelete) && application.GitOps.ProcessCode != 0 {
		// delete processdefine from ds
		client, dsSetting, err := c.getDsSettings(tenant)
		if err != nil {
			c.Logger.Error("error in getDsSettings: %v", err)
			return
		}
		if err := c.dsService.DeleteProcessDefine(client, dsSetting.Project, application.GitOps.ProcessCode); err != nil {
			c.Logger.Error(err)
			return
		}
		// delete processCode and cronId from database
		application.GitOps.ProcessCode = 0
		application.GitOps.CronId = 0
		c.svcCtx.Database.Model(&application).Updates(commontypes.Application{
			GitOps: application.GitOps,
		})
	}
}

func (c *CronService) AddGitOpsCron(application commontypes.Application, tenant string) (err error) {
	if application.GitOps.UseDS {
		return c.SaveDsJob(application, tenant)
	}
	var cronId cron.EntryID
	if cronId, err = c.svcCtx.Cron.AddFunc(application.GitOps.Cron, func() {
		if !c.CheckCron(application.AppName) {
			c.RemoveCron(application, tenant, false)
			return
		}
		if _, err = c.taskService.RunTask(&commontypes.RunTaskReq{
			AppName:     application.AppName,
			TaskType:    constant.TASK_TYPE_SYNCHRONIZE,
			TriggerType: constant.TASK_TRIGGER_TYPE_CRON,
			Waiting:     false,
		}, application, tenant); err != nil {
			c.Logger.Errorf("Failed to run autosync, application=%s: %w", application.AppName, err)
		}
	}); err != nil {
		return fmt.Errorf("failed to add cron for autosync, application=%s: %w", application.AppName, err)
	}
	application.GitOps.HostPortId = fmt.Sprintf("%s:%d", utils.GetLocalListener(c.svcCtx.Config.Port), cronId)
	c.svcCtx.Database.Model(&commontypes.Application{}).Where("app_name = ?", application.AppName).Updates(commontypes.Application{
		GitOps: application.GitOps,
	})
	c.svcCtx.CronIdMap[application.AppName] = types.CronData{
		CronId:    cronId,
		HostPort:  utils.GetLocalListener(c.svcCtx.Config.Port),
		Scheduler: application.GitOps.Cron,
	}
	c.Logger.Infof("Successfully add cron=%s for autosync, application=%s", application.GitOps.Cron, application.AppName)
	return
}

func (c *CronService) LoadCronJob() {
	// gitops cronjob
	var applications []commontypes.Application
	if err := c.svcCtx.Database.Model(&commontypes.Application{}).Where("git_ops like ?", `%自动%`).Find(&applications).Error; err != nil {
		c.Logger.Error(err)
		return
	}
	for _, application := range applications {
		if err := c.AddGitOpsCron(application, application.Tenant); err != nil {
			c.Logger.Error(err)
		}
	}
	c.svcCtx.Cron.Start()

	// task status check cronjob
	c.svcCtx.Cron.AddFunc("0 * * * *", func() {
		if c.svcCtx.LeaderElection.IsLeader() { // only leader do job
			c.taskService.CheckTaskStatus()
		}
	})
}

func (c *CronService) CheckCron(appName string) bool {
	var application commontypes.Application
	if err := c.svcCtx.Database.Model(&commontypes.Application{}).
		Where("app_name = ?", appName).
		Where("git_ops LIKE ?", "%自动%").
		First(&application).Error; err != nil {
		c.Logger.Infof("Check cron failed: application=%s not found or not '自动', cron will be deleted", appName)
		return false
	}
	if application.GitOps.HostPortId != fmt.Sprintf("%s:%d", c.svcCtx.CronIdMap[appName].HostPort, c.svcCtx.CronIdMap[appName].CronId) {
		c.Logger.Infof("Check cron failed: local cron hostPortId=%s:%d is not equal to database, cron will be deleted", c.svcCtx.CronIdMap[appName].HostPort, c.svcCtx.CronIdMap[appName].CronId)
		return false
	}
	return true
}

func (c *CronService) SaveDsJob(application commontypes.Application, tenant string) (err error) {
	// get dolphinscheduler settings
	client, dsSetting, e := c.getDsSettings(tenant)
	if e != nil {
		return e
	}
	var processCode int64
	if application.GitOps.ProcessCode == 0 {
		// get server_url settings
		setting := commontypes.Settings{}
		if err = c.svcCtx.Database.First(&setting, "setting_key = ?", "server_url").Error; err != nil {
			return fmt.Errorf("error finding server_url settings: %v", err)
		}
		serverUrl := setting.SettingValue

		// get processdefine template
		var tpl commontypes.YamlTemplate
		if err = c.svcCtx.Database.First(&tpl, "type = ?", "processdefine").Error; err != nil {
			return fmt.Errorf("error getting processdefine template: %v", err)
		}
		variables := map[string]interface{}{
			"TaskCode":    utils.GenerateRandomInt(15),
			"Appname":     application.AppName,
			"BearerToken": dsSetting.LizardcdJwtToken,
			"ServerUrl":   serverUrl,
		}
		var tmpl *template.Template
		if tmpl, err = template.New("processdefine").Parse(tpl.Content); err != nil {
			return fmt.Errorf("error in Parse processdefine template: %w", err)
		}
		var buf bytes.Buffer
		if err = tmpl.Execute(&buf, variables); err != nil {
			return fmt.Errorf("error in Execute processdefine template: %w", err)
		}
		jsonString := html.UnescapeString(buf.String())
		var jsonObj map[string]interface{}
		json.Unmarshal([]byte(jsonString), &jsonObj)

		// create processdefine
		if processCode, err = c.dsService.CreateProcessDefine(client, fmt.Sprintf("gitops/%s", application.AppName), dsSetting.Project, jsonObj); err != nil {
			return
		}
		// processdefine online
		if err = c.dsService.ReleaseProcessDefine(client, dsSetting.Project, processCode, "ONLINE"); err != nil {
			return
		}
		// create processdefine scheduler
		var cronId int64
		if cronId, err = c.dsService.CreateScheduler(client, dsSetting.Project, processCode, application.GitOps.Cron); err != nil {
			return
		}
		// scheduler online
		if err = c.dsService.ScheduleOnOffLine(client, dsSetting.Project, cronId, "online"); err != nil {
			return
		}
		// save projectCode and cronId to database
		application.GitOps.ProcessCode = processCode
		application.GitOps.CronId = cronId
		c.svcCtx.Database.Model(&application).Updates(commontypes.Application{
			GitOps: application.GitOps,
		})
	} else {
		// scheduler offline
		if err = c.dsService.ScheduleOnOffLine(client, dsSetting.Project, application.GitOps.CronId, "offline"); err != nil {
			return
		}
		// update scheduler
		if err = c.dsService.UpdateScheduler(client, dsSetting.Project, application.GitOps.CronId, application.GitOps.Cron); err != nil {
			return
		}
		// scheduler online
		if err = c.dsService.ScheduleOnOffLine(client, dsSetting.Project, application.GitOps.CronId, "online"); err != nil {
			return
		}
	}
	return
}

func (c *CronService) getDsSettings(tenant string) (client *req.Client, dsSetting *commontypes.DolphinschedulerSetting, err error) {
	var setting commontypes.Settings
	if err = c.svcCtx.Database.Model(&commontypes.Settings{}).
		Where("setting_key = ?", "dolphinscheduler").
		Where("tenant = ?", tenant).
		First(&setting).Error; err != nil {
		return nil, nil, fmt.Errorf("error finding dolphinscheduler settings: %v", err)
	}
	json.Unmarshal([]byte(setting.SettingValue), &dsSetting)
	client = utils.NewHttpClient(otel.Tracer("imroc/req")).
		SetBaseURL(dsSetting.BaseUrl).
		SetCommonHeader("token", dsSetting.Token)
	return
}
