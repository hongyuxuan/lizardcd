/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package helm

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var installVersion string
var cluster string
var namespace string
var valuesFile string

// installCmd represents the list command
var installCmd = &cobra.Command{
	Use:   "install [RELEASE_NAME] [REPO_NAME/CHART_NAME]",
	Short: "install helm charts",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		releaseName := args[0]
		chart := strings.Split(args[1], "/")
		if len(chart) != 2 {
			fmt.Printf("USAGE: %s helm install [RELEASE_NAME] [REPO_NAME/CHART_NAME] [flags]\n", common.GetExec())
			return
		}
		repoName := chart[0]
		chartName := chart[1]
		// get chart info
		var res *types.HelmSearchRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/helm/repo/%s?chart_name=%s", repoName, chartName)).SetResult(&res).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to search chart \"%s\" of helm repository \"%s\": %v", chartName, repoName, err)
		}
		chartInfo, ok := lo.Find(res.Data, func(c commontypes.ChartListResponse) bool {
			return c.ChartName == chartName
		})
		if !ok {
			common.PrintFatal("cannot find chart \"%s\" of helm repository \"%s\"", chartName, repoName)
		}
		params := map[string]interface{}{
			"repo_url":      chartInfo.RepoUrl,
			"chart_name":    chartName,
			"chart_version": installVersion,
			"release_name":  releaseName,
		}
		if valuesFile != "" {
			content, err := os.ReadFile(valuesFile)
			if err != nil {
				common.PrintFatal(err.Error())
				return
			}
			params["values"] = string(content)
		}
		if err := common.LizardServer.Post(fmt.Sprintf("/lizardcd/helm/cluster/%s/namespace/%s/charts/install", cluster, namespace)).SetBody(params).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to install chart \"%s\": %v", args[1], err)
		} else {
			common.PrintSuccess("successfully submit install chart \"%s-%s\"", args[1], installVersion)
		}
	},
}

func init() {
	installCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster (required)")
	installCmd.Flags().StringVar(&namespace, "namespace", "default", "kubernetes namespace")
	installCmd.Flags().StringVar(&installVersion, "version", "", "chart version (required)")
	installCmd.Flags().StringVarP(&valuesFile, "file", "f", "", "values.yaml")
	installCmd.MarkFlagRequired("cluster")
	installCmd.MarkFlagRequired("version")
}
