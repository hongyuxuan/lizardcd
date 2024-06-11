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
var valuesCmd = &cobra.Command{
	Use:   "values",
	Short: "show charts values.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		// get repo url
		var res *types.HelmRepoRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/db/helm_repositories?filter=name==%s", repo)).SetResult(&res).Do(context.Background()).Err; err != nil || res.Data.Total == 0 {
			common.PrintError("failed to find helm repository \"%s\": %v", repo, err)
			return
		}

		if err := common.LizardServer.Get("/lizardcd/helm/repo/charts/values").SetQueryParams(map[string]string{
			"repo_url":      res.Data.Results[0].URL,
			"chart_name":    chart,
			"chart_version": installVersion,
		}).SetOutput(os.Stdout).Do(context.Background()).Err; err != nil {
			common.PrintError("failed to show charts values.yaml: %v", err)
			return
		}
	},
}

func init() {
	valuesCmd.Flags().StringVar(&repo, "repo", "", "repo name (required)")
	valuesCmd.Flags().StringVar(&chart, "chart", "", "chart name (required)")
	valuesCmd.Flags().StringVar(&installVersion, "version", "", "chart version (required)")
	valuesCmd.MarkFlagRequired("repo")
	valuesCmd.MarkFlagRequired("chart")
	valuesCmd.MarkFlagRequired("version")
}
