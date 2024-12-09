/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package deployment

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/golang-module/carbon"
	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/svc"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/olekukonko/tablewriter"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/apps/v1"
)

var namespace string
var cluster string
var sortBy string
var order string

// listCmd represents the listdeployments command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List deployments in a namespace of a kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"name", "state", "lastupdate"})
		colors := tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor}
		table.SetHeaderColor(colors, colors, colors)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		var res *types.DeploymentRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/kubernetes/cluster/%s/namespace/%s/deployments", cluster, namespace)).SetResult(&res).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to get deployment list of cluster=%s, namespace=%s: %v", cluster, namespace, err)
		}

		var data [][]string
		for _, d := range res.Data {
			var lastUpdateTime string
			progress, ok := lo.Find(d.Status.Conditions, func(cond v1.DeploymentCondition) bool {
				return cond.Type == "Progressing"
			})
			if ok {
				lastUpdateTime = carbon.FromStdTime(progress.LastUpdateTime.Time).Format("Y-m-d H:i:s")
			}
			var state string = ""
			if d.Status.Replicas == 0 {
				state = "stopped"
			} else {
				if d.Status.ReadyReplicas >= d.Status.Replicas {
					state = "running"
				} else {
					state = "updating"
				}
			}
			state += fmt.Sprintf("(%d/%d)", d.Status.ReadyReplicas, d.Status.Replicas)
			row := []string{d.Name, state, lastUpdateTime}
			data = append(data, row)
		}
		svc.SortData(data, sortBy, order)

		for _, row := range data {
			var colors tablewriter.Colors
			if !common.Nocolor {
				if strings.HasPrefix(row[1], "running") {
					colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgGreenColor}
				} else if strings.HasPrefix(row[1], "updating") {
					colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgYellowColor}
				} else if strings.HasPrefix(row[1], "stopped") {
					colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor}
				}
			}
			table.Rich(row, []tablewriter.Colors{{}, colors, {}})
		}
		table.Render()
	},
}

func init() {
	listCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace of a kubernetes cluster")
	listCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster name")
	listCmd.Flags().StringVar(&sortBy, "sort-by", "name", "sort field of deployment, can be <lastUpdateTime|name>")
	listCmd.Flags().StringVar(&order, "order", "asc", "sort order of deployment, can be <asc|desc>")
	listCmd.MarkFlagRequired("cluster")
}
