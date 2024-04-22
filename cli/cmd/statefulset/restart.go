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

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart statefulset",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		if err := common.LizardServer.Patch(fmt.Sprintf("/kubernetes/cluster/%s/namespace/%s/statefulsets/%s/rollout", cluster, namespace, name)).Do(context.Background()).Err; err != nil {
			fmt.Printf("\033[0;31;40mrollout restart statefulset failed: %v\033[0m\n", err)
			return
		}
		fmt.Printf("\033[0;32;40mrollout restart statefulset success\033[0m\n")
	},
}

func init() {
	restartCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace of a kubernetes cluster")
	restartCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster name")
	restartCmd.Flags().StringVar(&name, "name", "", "statefulset name")
	restartCmd.MarkFlagRequired("cluster")
	restartCmd.MarkFlagRequired("name")
}
