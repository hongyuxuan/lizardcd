/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package deployment

import (
	"fmt"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

// deploymentCmd represents the deployment command
var DeploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "List/patch/rollout restart/scale of deployments",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Use \"%s deployment [command] --help\" for more information about a command.", common.GetExec())
	},
}

func init() {
	DeploymentCmd.AddCommand(listCmd)
	DeploymentCmd.AddCommand(showCmd)
	DeploymentCmd.AddCommand(setCmd)
	DeploymentCmd.AddCommand(restartCmd)
	DeploymentCmd.AddCommand(scaleCmd)
}
