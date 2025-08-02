/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/hongyuxuan/lizardcd/common/constant"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var appName string
var artifact string
var triggerType string
var labelStr string
var targetLabelStr string
var wait, markdown bool
var interval int64
var limitStr string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run task with app_name and artifact_url",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		var runRes *types.TaskExecuteRes
		labels := []string{}
		if labelStr != "" {
			labels = strings.Split(labelStr, ",")
		}
		targetLabels := []string{}
		if targetLabelStr != "" {
			targetLabels = strings.Split(targetLabelStr, ",")
		}
		limits := []string{}
		if limitStr != "" {
			limits = strings.Split(limitStr, ",")
		}

		// get application info, must use admin token
		adminToken := viper.GetString("lizardcd.auth.admin_token")
		var appRes types.ApplicationRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/db/application?filter=app_name==%s", appName)).
			SetBearerAuthToken(adminToken).
			SetSuccessResult(&appRes).
			Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to get application %s: %v", appName, err)
		}
		if appRes.Data.Total == 0 {
			common.PrintFatal("cannot find application %s", appName)
		}
		application := appRes.Data.Results[0]

		workloads := lo.Filter(application.Workload, func(item commontypes.Workload, _ int) bool {
			return item.Enable
		})

		if targetLabelStr != "" {
			workloads = lo.Filter(application.Workload, func(w commontypes.Workload, _ int) bool {
				return len(lo.Intersect(w.Labels, targetLabels)) > 0
			})
		}
		if limitStr != "" {
			workloads = lo.Filter(workloads, func(w commontypes.Workload, _ int) bool {
				return lo.Contains(limits, w.WorkloadName)
			})
		}
		if len(workloads) == 0 {
			common.PrintFatal("cannot find any targets to deploy")
		}
		for i := range workloads {
			workloads[i].ArtifactUrl = artifact
		}

		if err := common.LizardServer.Post("/lizardcd/task/run").
			SetBody(&commontypes.RunTaskReq{
				AppName:     appName,
				ArtifactUrl: artifact,
				TaskType:    "deploy",
				TriggerType: triggerType,
				Waiting:     false,
				Labels:      labels,
				Workloads:   workloads,
			}).
			SetSuccessResult(&runRes).
			Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to run task: %v", runRes.Data.Id, err)
		} else {
			link := fmt.Sprintf("%s/task/history/%s", viper.GetString("lizardcd.ui.url"), runRes.Data.Id)
			if markdown {
				common.PrintSuccess("#enableMarkdown\n任务执行详情：[%s](%s)", link, link)
			} else {
				common.PrintSuccess("任务执行详情：%s", link)
			}
		}

		// 获取结果
		if !wait {
			return
		}
		var res *types.TaskHistoryRes
		time.Sleep(time.Duration(interval) * time.Second) // wait a moment to fetch result
		for {
			if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/db/task_history/%s", runRes.Data.Id)).
				SetSuccessResult(&res).Do(context.Background()).Err; err != nil {
				common.PrintFatal(err.Error())
			}
			switch res.Data.Status {
			case constant.TASK_STATUS_RUNNING:
				common.PrintSuccess("app_name=%s task is running", appName)
			case constant.TASK_STATUS_FINISHED:
				if res.Data.Success.Bool {
					common.PrintSuccess("app_name=%s task finished successfully in %s", appName, res.Data.Expire)
				} else {
					common.PrintFatal("app_name=%s task finished failed: %v", appName, res.Data.ErrMessage)
				}
				return
			case constant.TASK_STATUS_INITIALIZE:
				if res.Data.Success.Valid {
					if !res.Data.Success.Bool {
						common.PrintFatal("app_name=%s task initialize failed: %v", appName, res.Data.ErrMessage)
						return
					}
				} else {
					common.PrintSuccess("app_name=%s task is initializing", appName)
				}
			case constant.TASK_STATUS_TERMINATED:
				if res.Data.Success.Bool {
					common.PrintSuccess("app_name=%s task terminated successfully in %s", appName, res.Data.Expire)
				} else {
					common.PrintFatal("app_name=%s task terminated failed: %v", appName, res.Data.ErrMessage)
				}
				return
			case constant.TASK_STATUS_WAITING:
				common.PrintSuccess("app_name=%s task is waiting", appName)
			}
			time.Sleep(time.Duration(interval) * time.Second)
		}
	},
}

func init() {
	runCmd.Flags().StringVar(&appName, "app-name", "", "Application name")
	runCmd.Flags().StringVar(&artifact, "artifact", "", "Artifact url to deploy for application")
	runCmd.Flags().StringVar(&triggerType, "trigger-type", "lizardcd-cli", "Trigger type for run task")
	runCmd.Flags().StringVar(&labelStr, "labels", "", "Labels for run task. eg, foo1:zoo1,foo2:zoo2")
	runCmd.Flags().StringVar(&targetLabelStr, "target-labels", "", "Host labels for deploy")
	runCmd.Flags().StringVar(&limitStr, "limits", "", "Hosts limit for deploy")
	runCmd.Flags().BoolVar(&wait, "wait", false, "If wait for finished")
	runCmd.Flags().BoolVar(&markdown, "markdown", false, "if print output with markdown")
	runCmd.Flags().Int64Var(&interval, "interval", 10, "Time interval to fetch run result")
	runCmd.MarkFlagRequired("app-name")
	runCmd.MarkFlagRequired("artifact")
}
