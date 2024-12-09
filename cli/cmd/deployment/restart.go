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

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart deployment",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		if err := common.LizardServer.Patch(fmt.Sprintf("/lizardcd/kubernetes/cluster/%s/namespace/%s/deployments/%s/rollout", cluster, namespace, name)).Do(context.Background()).Err; err != nil {
			common.PrintFatal("rollout restart deployment failed: %v", err)
		}
		common.PrintSuccess("rollout restart deployment:%s success", name)
	},
}

func init() {
	restartCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace of a kubernetes cluster")
	restartCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster name")
	restartCmd.Flags().StringVar(&name, "name", "", "deployment name")
	restartCmd.MarkFlagRequired("cluster")
	restartCmd.MarkFlagRequired("name")
}
