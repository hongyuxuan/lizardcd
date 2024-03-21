/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var AppVersion = "unknown"
var GoVersion = "unknown"
var BuildTime = "unknown"
var OsArch = "unknown"
var Author = "unknown"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Lizardcd-cli version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("App version: %s\n", AppVersion)
		fmt.Printf("Go version:  %s\n", GoVersion)
		fmt.Printf("Build Time:  %s\n", BuildTime)
		fmt.Printf("OS/Arch:     %s\n", OsArch)
		fmt.Printf("Author:      %s\n", Author)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
