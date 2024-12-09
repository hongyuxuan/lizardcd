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

// podCmd represents the pod command
var podCmd = &cobra.Command{
	Use:   "pod",
	Short: "List statefulset pods",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"podname", "podip", "nodename", "state", "message"})
		if !common.Nocolor {
			colors := tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor}
			table.SetHeaderColor(colors, colors, colors, colors, colors)
		}
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		data := svc.WritePodData(cluster, namespace, "statefulsets", name)
		for _, row := range data {
			colors := tablewriter.Colors{}
			if !common.Nocolor {
				if row[4] != "" {
					colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgRedColor}
				}
			}
			table.Rich(row, []tablewriter.Colors{{}, {}, {}, {}, colors})
		}
		table.Render()
	},
}

func init() {
	showCmd.AddCommand(podCmd)
	podCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace of a kubernetes cluster")
	podCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster name")
	podCmd.Flags().StringVar(&name, "name", "", "statefulset name")
	podCmd.MarkFlagRequired("cluster")
	podCmd.MarkFlagRequired("name")
}
