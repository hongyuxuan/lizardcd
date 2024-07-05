/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package statefulset

import (
	"os"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/svc"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// eventCmd represents the event command
var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "List events of statefulset pods",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"type", "reason", "age", "from", "message"})
		if !common.Nocolor {
			colors := tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor}
			table.SetHeaderColor(colors, colors, colors, colors, colors)
		}
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		data := svc.GetEvents(cluster, namespace, pod)
		for _, row := range data {
			colors := tablewriter.Colors{}
			if !common.Nocolor {
				if row[0] == "Warning" {
					colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgRedColor}
				}
			}
			table.Rich(row, []tablewriter.Colors{colors, colors, {}, {}, colors})
		}
		table.Render()
	},
}

func init() {
	showCmd.AddCommand(eventCmd)
	eventCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace of a kubernetes cluster")
	eventCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster name")
	eventCmd.Flags().StringVar(&pod, "pod", "", "pod name")
	eventCmd.MarkFlagRequired("cluster")
	eventCmd.MarkFlagRequired("pod")
}
