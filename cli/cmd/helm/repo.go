/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package helm

import (
	"fmt"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

// repoCmd represents the list command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "List/add/remove/update helm repositories",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Use \"%s helm repo [command] --help\" for more information about a command.\n", common.GetExec())
	},
}

func init() {
	repoCmd.AddCommand(repoListCmd)
	repoCmd.AddCommand(repoAddCmd)
	repoCmd.AddCommand(repoRemoveCmd)
	repoCmd.AddCommand(repoUpdateCmd)
}
