/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package config

import (
	common "github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverAddr string

// configCmd represents the config command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure lizardcd-cli config files",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()
		viper.Set("lizardcd.server.url", serverAddr)
		viper.WriteConfig()
	},
}

func init() {
	ConfigCmd.Flags().StringVarP(&serverAddr, "lizardcd.server.addr", "s", "http://localhost:5117", "lizardcd-server address")
}
