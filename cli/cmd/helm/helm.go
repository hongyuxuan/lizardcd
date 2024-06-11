/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package helm

import (
	"fmt"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

// HelmCmd represents the helm command
var HelmCmd = &cobra.Command{
	Use:   "helm",
	Short: "Helm management",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Use \"%s helm [command] --help\" for more information about a command.\n", common.GetExec())
	},
}

func init() {
	HelmCmd.AddCommand(repoCmd)
	HelmCmd.AddCommand(searchCmd)
	HelmCmd.AddCommand(installCmd)
	HelmCmd.AddCommand(uninstallCmd)
	HelmCmd.AddCommand(upgradeCmd)
	HelmCmd.AddCommand(rollbackCmd)
	HelmCmd.AddCommand(showCmd)
	HelmCmd.AddCommand(pullCmd)
	HelmCmd.AddCommand(releaseCmd)
}
