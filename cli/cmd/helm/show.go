/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package helm

import (
	"fmt"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

// showCmd represents the list command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "show charts values.yaml or readme.md",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Use \"%s helm show [command] --help\" for more information about a command.\n", common.GetExec())
	},
}

func init() {
	showCmd.AddCommand(valuesCmd)
	showCmd.AddCommand(readmeCmd)
}
