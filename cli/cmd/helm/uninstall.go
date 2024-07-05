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

// uninstallCmd represents the list command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall [RELEASE_NAME]",
	Short: "install helm charts",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		releaseName := args[0]
		// get chart info
		if err := common.LizardServer.Post(fmt.Sprintf("/lizardcd/helm/cluster/%s/namespace/%s/charts/uninstall?release_name=%s", cluster, namespace, releaseName)).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to uninstall release \"%s\": %v", releaseName, err)
		} else {
			common.PrintSuccess("successfully submit uninstall release \"%s\"", releaseName)
		}
	},
}

func init() {
	uninstallCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster (required)")
	uninstallCmd.Flags().StringVar(&namespace, "namespace", "default", "kubernetes namespace")
	uninstallCmd.MarkFlagRequired("cluster")
}
