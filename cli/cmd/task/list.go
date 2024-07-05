/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package task

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/golang-module/carbon"
	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var page int32
var limit int32
var time_from string
var time_till string

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"id", "app_name", "task_type", "trigger_type", "labels", "result", "status", "tenant", "init_at", "expire"})
		if !common.Nocolor {
			colors := tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor}
			table.SetHeaderColor(colors, colors, colors, colors, colors, colors, colors, colors, colors, colors)
		}
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		var res *types.TaskHistoriesRes
		url := fmt.Sprintf("/lizardcd/db/task_history?page=%d&size=%d&sort=init_at%%20desc", page, limit)
		if time_from != "" && time_till != "" {
			url += fmt.Sprintf("&range=init_at==%s,%s", time_from, time_till)
		}
		if err := common.LizardServer.Get(url).SetResult(&res).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to list tasks history: %v", err)
		}

		var data [][]string
		for _, d := range res.Data.Results {
			init_at := carbon.FromStdTime(d.InitAt.Time).Format("Y-m-d H:i:s")
			var result string
			if d.Success.Bool == true {
				result = "SUCCESS"
			} else {
				result = "FAIL"
			}
			row := []string{d.Id, d.AppName, d.TaskType, d.TriggerType, strings.Join(d.Labels, ", "), result, d.Status, d.Tenant, init_at, d.Expire}
			data = append(data, row)
		}

		for _, row := range data {
			var colors tablewriter.Colors
			if !common.Nocolor {
				if strings.HasPrefix(row[6], "running") {
					colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgYellowColor}
				} else {
					if strings.HasPrefix(row[5], "SUCCESS") {
						colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgGreenColor}
					} else if strings.HasPrefix(row[5], "FAIL") {
						colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgRedColor}
					}
				}
			}
			table.Rich(row, []tablewriter.Colors{{}, {}, {}, {}, {}, colors, colors, {}, {}, {}})
		}
		table.Render()
	},
}

func init() {
	listCmd.Flags().Int32Var(&page, "page", 1, "pages of listing application")
	listCmd.Flags().Int32Var(&limit, "limit", 10, "pageSize of listing application")
	listCmd.Flags().StringVar(&time_from, "time-from", "", "from time of listing tasks")
	listCmd.Flags().StringVar(&time_till, "time-till", "", "till time of listing tasks")
}
