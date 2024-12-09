/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package statefulset

import (
	"context"
	"fmt"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

var image string
var container string

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set image of a statefulset",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		if err := common.LizardServer.Patch(fmt.Sprintf("/lizardcd/kubernetes/cluster/%s/namespace/%s/statefulsets/%s?container=%s&image=%s", cluster, namespace, name, container, image)).Do(context.Background()).Err; err != nil {
			common.PrintFatal("set statefulset image failed: %v", err)
		}
		common.PrintSuccess("set statefulset:%s image:%s success", name, image)
	},
}

func init() {
	setCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace of a kubernetes cluster")
	setCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster name")
	setCmd.Flags().StringVar(&name, "name", "", "statefulset name")
	setCmd.Flags().StringVar(&container, "container", "", "statefulset container")
	setCmd.Flags().StringVar(&image, "image", "", "statefulset image")
	setCmd.MarkFlagRequired("cluster")
	setCmd.MarkFlagRequired("name")
	setCmd.MarkFlagRequired("container")
	setCmd.MarkFlagRequired("image")
}
