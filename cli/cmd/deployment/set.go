/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package deployment

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
	Short: "Set image of a deployment",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		if err := common.LizardServer.Patch(fmt.Sprintf("/lizardcd/kubernetes/cluster/%s/namespace/%s/deployments/%s?container=%s&image=%s", cluster, namespace, name, container, image)).Do(context.Background()).Err; err != nil {
			common.PrintFatal("set deployment image failed: %v", err)
		}
		common.PrintSuccess("set deployment:%s image:%s success", name, image)
	},
}

func init() {
	setCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace of a kubernetes cluster")
	setCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster name")
	setCmd.Flags().StringVar(&name, "name", "", "deployment name")
	setCmd.Flags().StringVar(&container, "container", "", "deployment container")
	setCmd.Flags().StringVar(&image, "image", "", "deployment image")
	setCmd.MarkFlagRequired("cluster")
	setCmd.MarkFlagRequired("name")
	setCmd.MarkFlagRequired("container")
	setCmd.MarkFlagRequired("image")
}
