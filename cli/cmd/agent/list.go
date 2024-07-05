/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package agent

import (
	"context"
	"os"

	common "github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get lizardcd agent list",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		var res *types.LizardAgentRes
		if err := common.LizardServer.Get("/lizardcd/server/services").SetResult(&res).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to get lizardcd agent list: %v\n", err)
		}

		var data [][]string
		for _, d := range res.Data {
			row := []string{d.ServiceName}
			data = append(data, row)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"lizardcd agent"})
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		if !common.Nocolor {
			table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor})
		}
		table.AppendBulk(data)
		table.Render()
	},
}
