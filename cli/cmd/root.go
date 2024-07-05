/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/hongyuxuan/lizardcd/cli/cmd/agent"
	"github.com/hongyuxuan/lizardcd/cli/cmd/application"
	"github.com/hongyuxuan/lizardcd/cli/cmd/apply"
	"github.com/hongyuxuan/lizardcd/cli/cmd/config"
	"github.com/hongyuxuan/lizardcd/cli/cmd/deployment"
	"github.com/hongyuxuan/lizardcd/cli/cmd/helm"
	"github.com/hongyuxuan/lizardcd/cli/cmd/login"
	"github.com/hongyuxuan/lizardcd/cli/cmd/statefulset"
	"github.com/hongyuxuan/lizardcd/cli/cmd/task"
	common "github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   common.GetExec(),
	Short: "A command-line tools for lizardcd, which is a cloud native continuous delivery for kubernetes",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&common.ConfigFile, "config", "c", "~/.lizardcd-cli.yaml", "config file")
	rootCmd.PersistentFlags().StringVarP(&common.LogLevel, "log.level", "l", "info", "log level")
	rootCmd.PersistentFlags().BoolVar(&common.Nocolor, "nocolor", false, "if ouput without color")

	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(agent.AgentCmd)
	rootCmd.AddCommand(deployment.DeploymentCmd)
	rootCmd.AddCommand(statefulset.StatefulsetCmd)
	rootCmd.AddCommand(apply.ApplyCmd)
	rootCmd.AddCommand(login.LoginCmd)
	rootCmd.AddCommand(application.ApplicationCmd)
	rootCmd.AddCommand(task.TaskCmd)
	rootCmd.AddCommand(helm.HelmCmd)
}
