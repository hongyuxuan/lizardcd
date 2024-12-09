/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package statefulset

import (
	"os"
	"strings"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/svc"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var pod string

// containerCmd represents the container command
var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "List statefulset containers",

	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"container", "status", "state", "restart_count"})
		if !common.Nocolor {
			colors := tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor}
			table.SetHeaderColor(colors, colors, colors, colors)
		}
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoWrapText(false)

		p := svc.GetPod(cluster, namespace, "statefulsets", name, pod)
		var data [][]string
		data = append(data, svc.WriteContainerData(p.Status.ContainerStatuses, p.Spec.Containers, p.Status.Conditions, false)...)
		data = append(data, svc.WriteContainerData(p.Status.InitContainerStatuses, p.Spec.InitContainers, p.Status.Conditions, true)...)

		for _, row := range data {
			var colors tablewriter.Colors
			var restartColor tablewriter.Colors
			if !common.Nocolor {
				restartColor = tablewriter.Colors{tablewriter.Normal, tablewriter.FgCyanColor}
				if strings.HasPrefix(row[1], "running") {
					colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgGreenColor}
				} else if strings.HasPrefix(row[1], "waiting") {
					colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgYellowColor}
				} else if strings.HasPrefix(row[1], "terminated") {
					colors = tablewriter.Colors{tablewriter.Normal, tablewriter.FgHiBlackColor}
				}
			}
			table.Rich(row, []tablewriter.Colors{{}, colors, colors, restartColor})
		}
		table.Render()
	},
}

func init() {
	showCmd.AddCommand(containerCmd)
	containerCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace of a kubernetes cluster")
	containerCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster name")
	containerCmd.Flags().StringVar(&name, "name", "", "statefulset name")
	containerCmd.Flags().StringVar(&pod, "pod", "", "pod name")
	containerCmd.MarkFlagRequired("cluster")
	containerCmd.MarkFlagRequired("name")
	containerCmd.MarkFlagRequired("pod")
}
