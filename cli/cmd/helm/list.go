/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package helm

import (
	"context"
	"os"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// repoListCmd represents the list command
var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "List helm repositories",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"name", "url", "tenant"})
		if !common.Nocolor {
			colors := tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor}
			table.SetHeaderColor(colors, colors, colors)
		}
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		var res *types.HelmRepoRes
		if err := common.LizardServer.Get("/lizardcd/db/helm_repositories").SetResult(&res).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to list helm repositories: %v", err)
		}

		var data [][]string
		for _, d := range res.Data.Results {
			row := []string{d.Name, d.URL, d.Tenant}
			data = append(data, row)
		}
		table.AppendBulk(data)
		table.Render()
	},
}

func init() {
}
