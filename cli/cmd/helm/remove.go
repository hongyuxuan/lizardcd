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
var repoRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove helm repositories",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		if err := common.LizardServer.Delete(fmt.Sprintf("/lizardcd/helm/repo/%s", name)).SetBody(map[string]string{
			"name":     name,
			"url":      url,
			"username": username,
			"password": password,
		}).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to add helm repository: %v", err)
		}
		common.PrintSuccess("successfully remove helm repository \"%s\"", name)
	},
}

func init() {
	repoRemoveCmd.Flags().StringVar(&name, "name", "", "helm repository name (required)")
	repoRemoveCmd.MarkFlagRequired("name")
}
