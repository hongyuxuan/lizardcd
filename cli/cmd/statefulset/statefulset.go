/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package statefulset

import (
	"fmt"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

// statefulsetCmd represents the statefulset command
var StatefulsetCmd = &cobra.Command{
	Use:   "statefulset",
	Short: "List/patch/rollout restart/scale of deployments",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Use \"%s statefulset [command] --help\" for more information about a command.", common.GetExec())
	},
}

func init() {
	StatefulsetCmd.AddCommand(listCmd)
	StatefulsetCmd.AddCommand(showCmd)
}
