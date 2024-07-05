/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package helm

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var version bool

// searchCmd represents the list command
var searchCmd = &cobra.Command{
	Use:   "search [REPO_NAME/CHART_NAME]",
	Short: "search chart of helm repositories",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		chart := strings.Split(args[0], "/")
		if len(chart) != 2 {
			fmt.Printf("USAGE: %s helm search [REPO_NAME/CHART_NAME] [flags]\n", common.GetExec())
		}
		repoName := chart[0]
		chartName := chart[1]

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"chart_name", "chart_version", "app_version", "description"})
		if !common.Nocolor {
			colors := tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor}
			table.SetHeaderColor(colors, colors, colors, colors)
		}
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		var res *types.HelmSearchRes
		if !version {
			if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/helm/repo/%s?chart_name=%s", repoName, chartName)).SetResult(&res).Do(context.Background()).Err; err != nil {
				common.PrintFatal("failed to search chart \"%s\" of helm repository \"%s\": %v", chartName, repoName, err)
			}
		} else {
			if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/helm/repo/%s/%s", repoName, chartName)).SetResult(&res).Do(context.Background()).Err; err != nil {
				common.PrintFatal("failed to search chart \"%s\" of helm repository \"%s\": %v", chartName, repoName, err)
			}
		}

		var data [][]string
		for _, d := range res.Data {
			var desc string
			if len(d.Description) > 80 {
				desc = d.Description[0:80] + "..."
			}
			row := []string{repoName + "/" + d.ChartName, d.ChartVersion, d.AppVersion, desc}
			data = append(data, row)
		}
		table.AppendBulk(data)
		table.Render()
	},
}

func init() {
	searchCmd.Flags().BoolVar(&version, "version", false, "if list chart version")
}
