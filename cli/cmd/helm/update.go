/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package helm

import (
	"context"
	"fmt"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

// repoListCmd represents the list command
var repoUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update helm repositories",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		if err := common.LizardServer.Post(fmt.Sprintf("/lizardcd/helm/cluster/%s/namespace/%s/repo/update", cluster, namespace)).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to update helm repository of cluster=%s namespace=%s: %v", cluster, namespace, err)
		}
		common.PrintSuccess("successfully update helm repository of cluster=%s namespace=%s", cluster, namespace)
	},
}

func init() {
	repoUpdateCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster (required)")
	repoUpdateCmd.Flags().StringVar(&namespace, "namespace", "default", "kubernetes namespace")
	repoUpdateCmd.MarkFlagRequired("cluster")
}
