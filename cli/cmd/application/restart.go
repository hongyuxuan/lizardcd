/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package application

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart application of all of it's workloads",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		var res *types.ApplicationRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/db/application?page=1&size=1&filter=app_name==%s", url.QueryEscape(application))).SetResult(&res).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to get application \"%s\": %v", application, err)
		}
		if res.Data.Total == 0 {
			common.PrintFatal("cannot find application \"%s\"", application)
		}

		app := res.Data.Results[0]
		if err := common.LizardServer.Post("/lizardcd/task/run").SetBody(map[string]interface{}{
			"app_name":     application,
			"task_type":    "rollout",
			"trigger_type": "lizardcd-cli",
			"workloads": lo.Map(app.Workload, func(w commontypes.WorkLoad, _ int) map[string]interface{} {
				return map[string]interface{}{
					"cluster":       w.Cluster,
					"namespace":     w.Namespace,
					"workload_type": w.WorkloadType,
					"workload_name": w.WorkloadName,
				}
			}),
		}).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to rollout restart application \"%s\": %v", application, err)
		} else {
			common.PrintSuccess("successfully submit rollout restart task, use \"%s task list\" to see results", common.GetExec())
		}
	},
}

func init() {
	restartCmd.Flags().StringVar(&application, "name", "", "application name")
	restartCmd.MarkFlagRequired("name")
}
