/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package statefulset

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
	Short: "Set statefulset replicas",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		if err := common.LizardServer.Patch(fmt.Sprintf("/kubernetes/cluster/%s/namespace/%s/statefulsets/scale", cluster, namespace)).SetBody(&types.ScaleReq{
			Workloads: []types.Workloads{
				{
					Name:     name,
					Replicas: replicas,
				},
			},
		}).Do(context.Background()).Err; err != nil {
			fmt.Printf("\033[0;31;40mset statefulset replicas failed: %v\033[0m\n", err)
			return
		}
		fmt.Printf("\033[0;32;40mset statefulset replicas success\033[0m\n")
	},
}

func init() {
	scaleCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace of a kubernetes cluster")
	scaleCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster name")
	scaleCmd.Flags().StringVar(&name, "name", "", "statefulset name")
	scaleCmd.Flags().Int32Var(&replicas, "replicas", 1, "statefulset replicas")
	scaleCmd.MarkFlagRequired("cluster")
	scaleCmd.MarkFlagRequired("name")
	scaleCmd.MarkFlagRequired("replicas")
}
