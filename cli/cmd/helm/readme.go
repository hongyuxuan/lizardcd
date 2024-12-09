/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package helm

import (
	"context"
	"fmt"
	"os"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/spf13/cobra"
)

// repoListCmd represents the list command
var readmeCmd = &cobra.Command{
	Use:   "readme",
	Short: "show charts readme.md",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		// get repo url
		var res *types.HelmRepoRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/db/helm_repositories?filter=name==%s", repo)).SetResult(&res).Do(context.Background()).Err; err != nil || res.Data.Total == 0 {
			common.PrintFatal("failed to find helm repository \"%s\": %v", repo, err)
		}

		if err := common.LizardServer.Get("/lizardcd/helm/repo/charts/readme").SetQueryParams(map[string]string{
			"repo_url":      res.Data.Results[0].URL,
			"chart_name":    chart,
			"chart_version": installVersion,
		}).SetOutput(os.Stdout).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to show charts readme.md: %v", err)
		}
	},
}

func init() {
	readmeCmd.Flags().StringVar(&repo, "repo", "", "repo name (required)")
	readmeCmd.Flags().StringVar(&chart, "chart", "", "chart name (required)")
	readmeCmd.Flags().StringVar(&installVersion, "version", "", "chart version (required)")
	readmeCmd.MarkFlagRequired("repo")
	readmeCmd.MarkFlagRequired("chart")
	readmeCmd.MarkFlagRequired("version")
}
