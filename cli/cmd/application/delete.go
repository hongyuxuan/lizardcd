/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package application

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/spf13/cobra"
)

// deleteCmd represents the restart command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete an application",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		var res *types.ApplicationRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/db/application?page=1&size=1&filter=app_name==%s", url.QueryEscape(application))).SetResult(&res).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to get application \"%s\": %v", application, err)
		}
		if res.Data.Total == 0 {
			common.PrintFatal("cannot find application \"%s\"", application)
		}

		app := res.Data.Results[0]
		if err := common.LizardServer.Delete(fmt.Sprintf("/lizardcd/db/application/%d", app.Id)).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to delete application \"%s\": %v", application, err)
		} else {
			common.PrintSuccess("successfully delete application \"%s\"", application)
		}
	},
}

func init() {
	deleteCmd.Flags().StringVar(&application, "name", "", "application name (required)")
	deleteCmd.MarkFlagRequired("name")
}
