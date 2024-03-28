/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package agent

import (
	"context"
	"os"

	common "github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
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
		if err := common.LizardServer.Get("/consul/services").SetResult(&res).Do(context.Background()).Err; err != nil {
			utils.Log.Fatalf("failed to get lizardcd agent list: %v", err)
		}

		var data [][]string
		for _, d := range res.Data {
			row := []string{d.ServiceName}
			data = append(data, row)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"lizardcd agent"})
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor})
		table.AppendBulk(data)
		table.Render()
	},
}
