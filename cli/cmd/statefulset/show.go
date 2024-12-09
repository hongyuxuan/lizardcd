/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package statefulset

import (
	"fmt"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

var name string

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "List statefulset pods or containers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Use \"%s statefulset show [command] --help\" for more information about a command.", common.GetExec())
	},
}
