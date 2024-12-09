/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package helm

import (
	"fmt"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

// releaseCmd represents the list command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "List/reinstall/rollback releases installed by helm charts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Use \"%s helm release [command] --help\" for more information about a command.\n", common.GetExec())
	},
}

func init() {
	releaseCmd.AddCommand(releaseListCmd)
}
