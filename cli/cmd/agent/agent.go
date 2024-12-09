/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package agent

import (
	"fmt"

	common "github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var AgentCmd = &cobra.Command{
	Use:   "agent",
	Short: "List lizardcd agent and/or other commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Use \"%s agent [command] --help\" for more information about a command.", common.GetExec())
	},
}

func init() {
	AgentCmd.AddCommand(ListCmd)
}
