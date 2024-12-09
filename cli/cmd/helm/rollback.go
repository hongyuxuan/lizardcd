/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package helm

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

// rollbackCmd represents the list command
var rollbackCmd = &cobra.Command{
	Use:   "rollback [RELEASE_NAME] [REVISION]",
	Short: "install helm charts",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		releaseName := args[0]
		revision, _ := strconv.ParseInt(args[1], 10, 32)
		if err := common.LizardServer.Post(fmt.Sprintf("/lizardcd/helm/cluster/%s/namespace/%s/release/rollback", cluster, namespace)).SetBody(map[string]interface{}{
			"release_name": releaseName,
			"revision":     revision,
		}).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to rollback release \"%s\": %v", releaseName, err)
		} else {
			common.PrintSuccess("successfully submit rollback release \"%s\"", releaseName)
		}
	},
}

func init() {
	rollbackCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster (required)")
	rollbackCmd.Flags().StringVar(&namespace, "namespace", "default", "kubernetes namespace")
	rollbackCmd.MarkFlagRequired("cluster")
}
