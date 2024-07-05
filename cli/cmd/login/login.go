/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package login

import (
	"context"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/common/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var username string
var password string
var save bool

// LoginCmd represents the login command
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to lizardcd-server",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()
		var res *types.LoginRes
		if err := common.LizardServer.Post("/lizardcd/auth/login").SetBody(map[string]string{
			"username": username,
			"password": password,
		}).SetResult(&res).Do(context.Background()).Err; err != nil {
			common.PrintFatal(err.Error())
		}
		viper.Set("lizardcd.auth.access_token", res.AccessToken)
		if save {
			viper.Set("lizardcd.auth.username", username)
			viper.Set("lizardcd.auth.password", password)
			viper.WriteConfig()
		}
	},
}

func init() {
	LoginCmd.Flags().StringVarP(&username, "username", "u", "", "lizardcd-server username")
	LoginCmd.Flags().StringVarP(&password, "password", "p", "", "lizardcd-server password")
	LoginCmd.Flags().BoolVar(&save, "save", false, "if save lizardcd-server password")
	LoginCmd.MarkFlagRequired("username")
	LoginCmd.MarkFlagRequired("password")
}
