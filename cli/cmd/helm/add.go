/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package helm

import (
	"context"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

var name string
var url string
var username string
var password string

// repoListCmd represents the list command
var repoAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add helm repositories",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		if err := common.LizardServer.Post("/lizardcd/helm/repo").SetBody(map[string]string{
			"name":     name,
			"url":      url,
			"username": username,
			"password": password,
		}).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to add helm repository: %v", err)
		}
		common.PrintSuccess("successfully add helm repository \"%s\"", name)
	},
}

func init() {
	repoAddCmd.Flags().StringVar(&name, "name", "", "helm repository name (required)")
	repoAddCmd.Flags().StringVar(&url, "url", "", "helm repository url (required)")
	repoAddCmd.Flags().StringVar(&username, "username", "", "chart version")
	repoAddCmd.Flags().StringVar(&password, "password", "", "values.yaml")
	repoAddCmd.MarkFlagRequired("name")
	repoAddCmd.MarkFlagRequired("url")
}
