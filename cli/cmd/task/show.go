/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package task

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/golang-module/carbon"
	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/olekukonko/tablewriter"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var id string

type podStatus struct {
	PodName string `json:"pod_name"`
	Ready   string `json:"ready"`
}

// showCmd represents the list command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show task history detail",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"cluster", "namespace", "deployment_type", "deployment_name", "container_name", "image", "status", "err_message", "update_at"})
		if !common.Nocolor {
			colors := tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor}
			table.SetHeaderColor(colors, colors, colors, colors, colors, colors, colors, colors, colors)
		}
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		var res *types.TaskHistoryRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/db/task_history/%s", id)).SetResult(&res).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to get task history: %v", err)
		}

		var data [][]string
		for _, w := range res.Data.TaskHistoryWorkloads {
			update_at := carbon.FromStdTime(w.UpdateAt).Format("Y-m-d H:i:s")
			var podstatus []podStatus
			json.Unmarshal([]byte(w.Status), &podstatus)
			status := lo.Map(podstatus, func(p podStatus, _ int) string {
				return p.PodName + ":" + p.Ready
			})
			row := []string{w.Workload.Cluster, w.Workload.Namespace, w.Workload.WorkloadType, w.Workload.WorkloadName, w.Workload.ContainerName, w.Workload.ArtifactUrl, strings.Join(status, ","), w.ErrMessage, update_at}
			data = append(data, row)
		}

		for _, row := range data {
			var colors tablewriter.Colors
			if !common.Nocolor {
				if row[8] != "" {
					colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgRedColor}
				}
			}
			table.Rich(row, []tablewriter.Colors{{}, {}, {}, {}, {}, {}, {}, colors, {}})
		}
		table.Render()
	},
}

func init() {
	showCmd.Flags().StringVar(&id, "id", "", "task id")
	showCmd.MarkFlagRequired("id")
}
