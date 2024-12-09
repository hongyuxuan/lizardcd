/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package helm

import (
	"context"
	"fmt"
	"strings"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/spf13/cobra"
)

var repo string
var chart string

// pullCmd represents the list command
var pullCmd = &cobra.Command{
	Use:   "pull [REPO_NAME/CHART_NAME]",
	Short: "download charts",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		chart := strings.Split(args[0], "/")
		if len(chart) != 2 {
			fmt.Printf("USAGE: %s helm pull [REPO_NAME/CHART_NAME] [flags]\n", common.GetExec())
			return
		}
		repoName := chart[0]
		chartName := chart[1]

		// get repo url
		var res *types.HelmRepoRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/db/helm_repositories?filter=name==%s", repoName)).SetResult(&res).Do(context.Background()).Err; err != nil || res.Data.Total == 0 {
			common.PrintFatal("failed to find helm repository \"%s\": %v", repoName, err)
		}

		if err := common.LizardServer.Get("/lizardcd/helm/repo/charts/download").SetQueryParams(map[string]string{
			"repo_url":      res.Data.Results[0].URL,
			"chart_name":    chartName,
			"chart_version": installVersion,
		}).SetOutputFile(fmt.Sprintf("%s-%s.tgz", chartName, installVersion)).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to download charts: %v", err)
			return
		}
	},
}

func init() {
	pullCmd.Flags().StringVar(&installVersion, "version", "", "chart version (required)")
	pullCmd.MarkFlagRequired("version")
}
