/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package application

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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List applications",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"app_name", "tags", "tenant", "update_at"})
		if !common.Nocolor {
			colors := tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor}
			table.SetHeaderColor(colors, colors, colors, colors)
		}
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		var res *types.ApplicationRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/db/application?page=%d&size=%d&sort=update_at%%20desc", page, limit)).SetResult(&res).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to list application: %v", err)
		}

		var data [][]string
		for _, d := range res.Data.Results {
			update_at := carbon.FromStdTime(d.UpdateAt).Format("Y-m-d H:i:s")
			tags := d.Tags
			if tags == nil {
				tags = []string{}
			}
			row := []string{d.AppName, strings.Join(tags, ", "), d.Tenant, update_at}
			data = append(data, row)
		}

		table.AppendBulk(data)
		table.Render()
	},
}

func init() {
	listCmd.Flags().Int32Var(&page, "page", 1, "pages of listing application")
	listCmd.Flags().Int32Var(&limit, "limit", 10, "pageSize of listing application")
}
