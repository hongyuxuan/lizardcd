/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package task

import (
	"context"
	"fmt"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	"github.com/spf13/cobra"
)

// showCmd represents the list command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "execute task by id",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		var res *types.TaskHistoryRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/db/task_history/%s", id)).SetResult(&res).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to get task history with id=%s: %v", id, err)
		}

		var executeRes *types.TaskExecuteRes
		if err := common.LizardServer.Post(fmt.Sprintf("/lizardcd/task/execute/%s", id)).SetResult(&executeRes).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to execute task with id=%s: %v", id, err)
		} else {
			common.PrintSuccess("successful submit task, id=%s", id)
		}
	},
}

func init() {
	executeCmd.Flags().StringVar(&id, "id", "", "task id (required)")
	executeCmd.MarkFlagRequired("id")
}
