/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package task

import (
	"fmt"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

// taskCmd represents the task command
var TaskCmd = &cobra.Command{
	Use:   "task",
	Short: "List/show/execute/run task history",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Use \"%s task [command] --help\" for more information about a command.", common.GetExec())
	},
}

func init() {
	TaskCmd.AddCommand(listCmd)
	TaskCmd.AddCommand(showCmd)
	TaskCmd.AddCommand(executeCmd)
	TaskCmd.AddCommand(runCmd)
}
