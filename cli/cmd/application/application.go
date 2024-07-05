/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package application

import (
	"fmt"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

// applicationCmd represents the application command
var ApplicationCmd = &cobra.Command{
	Use:   "application",
	Short: "List/deploy/restart/delete applications",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Use \"%s application [command] --help\" for more information about a command.", common.GetExec())
	},
}

func init() {
	ApplicationCmd.AddCommand(listCmd)
	ApplicationCmd.AddCommand(deployCmd)
	ApplicationCmd.AddCommand(restartCmd)
}
