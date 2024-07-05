/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package helm

import (
	"context"
	"fmt"
	"os"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// releaseListCmd represents the list command
var releaseListCmd = &cobra.Command{
	Use:   "list",
	Short: "List releases installed by helm charts",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"name", "chart_version", "app_version", "revision", "status", "update_at"})
		if !common.Nocolor {
			colors := tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor}
			table.SetHeaderColor(colors, colors, colors, colors, colors, colors)
		}
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		var res *types.HelmReleaseRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/helm/cluster/%s/namespace/%s/releases", cluster, namespace)).SetResult(&res).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to list releases: %v", err)
			return
		}

		var data [][]string
		for _, d := range res.Data {
			row := []string{d.Name, d.Chart, d.AppVersion, d.Revision, d.Status, d.Updated}
			data = append(data, row)
		}
		for _, row := range data {
			var colors tablewriter.Colors
			if row[4] == "deployed" || row[4] == "uninstalled" {
				colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgGreenColor}
			} else if row[4] == "failed" {
				colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgRedColor}
			} else {
				colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgYellowColor}
			}
			table.Rich(row, []tablewriter.Colors{{}, {}, {}, {}, colors, {}})
		}
		table.Render()
	},
}

func init() {
	releaseListCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster (required)")
	releaseListCmd.Flags().StringVar(&namespace, "namespace", "default", "kubernetes namespace")
	releaseListCmd.MarkFlagRequired("cluster")
}
