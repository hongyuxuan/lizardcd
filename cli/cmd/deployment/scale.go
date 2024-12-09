/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package deployment

import (
	"context"
	"fmt"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/spf13/cobra"
)

var replicas int32

// scaleCmd represents the scale command
var scaleCmd = &cobra.Command{
	Use:   "scale",
	Short: "Set deployment replicas",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		if err := common.LizardServer.Patch(fmt.Sprintf("/lizardcd/kubernetes/cluster/%s/namespace/%s/deployments/scale", cluster, namespace)).SetBody(&types.ScaleReq{
			Workloads: []types.Workloads{
				{
					Name:     name,
					Replicas: replicas,
				},
			},
		}).Do(context.Background()).Err; err != nil {
			common.PrintFatal("set deployment replicas failed: %v", err)
		}
		common.PrintSuccess("set deployment:%s replicas:%d success", name, replicas)
	},
}

func init() {
	scaleCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace of a kubernetes cluster")
	scaleCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster name")
	scaleCmd.Flags().StringVar(&name, "name", "", "deployment name")
	scaleCmd.Flags().Int32Var(&replicas, "replicas", 1, "deployment replicas")
	scaleCmd.MarkFlagRequired("cluster")
	scaleCmd.MarkFlagRequired("name")
	scaleCmd.MarkFlagRequired("replicas")
}
