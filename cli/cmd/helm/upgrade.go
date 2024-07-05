/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package helm

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/spf13/cobra"
)

// upgradeCmd represents the list command
var upgradeCmd = &cobra.Command{
	Use:   "upgrade [RELEASE_NAME] [REPO_NAME/CHART_NAME]",
	Short: "upgrade helm charts",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		releaseName := args[0]
		chart := strings.Split(args[1], "/")
		if len(chart) != 2 {
			fmt.Printf("USAGE: %s helm upgrade [RELEASE_NAME] [REPO_NAME/CHART_NAME] [flags]\n", common.GetExec())
			return
		}
		repoName := chart[0]
		chartName := chart[1]
		// get repo url
		var res *types.HelmRepoRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/db/helm_repositories?filter=name==%s", repoName)).SetResult(&res).Do(context.Background()).Err; err != nil || res.Data.Total == 0 {
			common.PrintFatal("failed to find helm repository \"%s\": %v", repoName, err)
		}
		// get release info
		var release *types.HelmReleaseRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/helm/cluster/%s/namespace/%s/releases?release_name=%s", cluster, namespace, releaseName)).SetResult(&release).Do(context.Background()).Err; err != nil || len(release.Data) == 0 {
			common.PrintFatal("failed to find helm release \"%s\": %v", releaseName, err)
		}
		revision, _ := strconv.ParseInt(release.Data[0].Revision, 10, 32)
		params := map[string]interface{}{
			"repo_url":      res.Data.Results[0].URL,
			"chart_name":    chartName,
			"chart_version": release.Data[0].ChartVersion,
			"revision":      revision,
			"release_name":  releaseName,
		}
		if valuesFile != "" {
			content, err := os.ReadFile(valuesFile)
			if err != nil {
				common.PrintFatal(err.Error())
			}
			params["values"] = string(content)
		}
		if installVersion != "" {
			params["chart_version"] = installVersion
		}
		if err := common.LizardServer.Post(fmt.Sprintf("/lizardcd/helm/cluster/%s/namespace/%s/charts/upgrade", cluster, namespace)).SetBody(params).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to upgrade chart \"%s\": %v", args[1], err)
		} else {
			common.PrintSuccess("successfully submit upgrade chart \"%s-%s\"", args[1], installVersion)
		}
	},
}

func init() {
	upgradeCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster (required)")
	upgradeCmd.Flags().StringVar(&namespace, "namespace", "default", "kubernetes namespace")
	upgradeCmd.Flags().StringVar(&installVersion, "version", "", "chart version")
	upgradeCmd.Flags().StringVarP(&valuesFile, "file", "f", "", "values.yaml")
	upgradeCmd.MarkFlagRequired("cluster")
}
