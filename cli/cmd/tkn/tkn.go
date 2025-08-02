/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package tkn

import (
	"fmt"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

// taskCmd represents the task command
var TknCmd = &cobra.Command{
	Use:   "tkn",
	Short: "Start a tekton pipelinerun",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Use \"%s tkn [command] --help\" for more information about a command.", common.GetExec())
	},
}

func init() {
	TknCmd.AddCommand(startCmd)
}
