/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package database

import (
	"fmt"

	common "github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/spf13/cobra"
)

// databaseCmd represents the database command
var DatabaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database command: initdb,migrate",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Use \"%s database [command] --help\" for more information about a command.", common.GetExec())
	},
}

func init() {
	DatabaseCmd.AddCommand(MigrateCmd)
}
