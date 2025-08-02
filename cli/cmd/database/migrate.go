/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package database

import (
	"strings"

	"github.com/hongyuxuan/lizardcd/cli/common"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var dbfile, dbHost, username, password, db, tables string
var dbPort int

var srcdb *gorm.DB
var destdb *gorm.DB

var offset = 0
var size = 100

var funcMap map[string]func()

// migrateCmd represents the migrate command
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database from sqlite to tidb/mysql/mariadb",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		srcdb = utils.NewSQLite(dbfile, common.LogLevel)
		destdb = utils.NewTidb(dbHost, username, password, db, common.LogLevel, dbPort)

		funcMap = map[string]func(){
			"agent":                 migrate_agent,
			"application":           migrate_application,
			"application_faas":      migrate_application_faas,
			"application_resource":  migrate_application_resource,
			"ci_trigger":            migrate_ci_trigger,
			"git_repository":        migrate_git_repository,
			"helm_repository":       migrate_helm_repository,
			"image_repository":      migrate_image_repository,
			"oauth2":                migrate_oauth2,
			"settings":              migrate_settings,
			"task_history":          migrate_task_history,
			"task_history_workload": migrate_task_history_workload,
			"tenant":                migrate_tenant,
			"token":                 migrate_token,
			"user":                  migrate_user,
			"yaml_template":         migrate_yaml_template,
		}
		for _, table := range strings.Split(tables, ",") {
			if f, exists := funcMap[table]; exists {
				f()
			}
		}
	},
}

func init() {
	MigrateCmd.Flags().StringVarP(&dbfile, "db-file", "d", "./lizardcd.db", "database file, only for sqlite")
	MigrateCmd.Flags().StringVarP(&dbHost, "db-host", "H", "localhost", "database host")
	MigrateCmd.Flags().IntVarP(&dbPort, "db-port", "P", 3306, "database port")
	MigrateCmd.Flags().StringVarP(&username, "username", "u", "lizardcd", "database username")
	MigrateCmd.Flags().StringVarP(&password, "password", "p", "", "database password")
	MigrateCmd.Flags().StringVar(&db, "db", "lizardcd", "database db")
	MigrateCmd.Flags().StringVarP(&tables, "tables", "t", "", "database tables")
}

func migrate_agent() {
	utils.Log.Infof("========== Begin to migrate agent ==========")
	var agents []commontypes.Agent
	offset = 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&agents).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(agents) == 0 {
			break
		}
		result := destdb.Create(&agents)
		if result.Error != nil {
			utils.Log.Warnf("Batch migrate agents failed - %v", result.Error)
			offset += size
			continue
		}
		utils.Log.Infof("Batch migrate agents success, rows affected=%d", result.RowsAffected)
		offset += size
	}
}

func migrate_application() {
	// application
	utils.Log.Infof("========== Begin to migrate application ==========")
	var applications []commontypes.Application
	offset := 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&applications).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(applications) == 0 {
			break
		}
		for _, item := range applications {
			if item.AutoBuild == nil {
				item.AutoBuild = &commontypes.AutoBuild{}
			}
			if item.GitOps == nil {
				item.GitOps = &commontypes.GitOps{}
			}
			if err := destdb.Save(&item).Error; err != nil {
				utils.Log.Warnf("id=%d app_name=%s - %v", item.Id, item.AppName, err)
				continue
			}
			utils.Log.Infof("id=%d app_name=%s migrate success", item.Id, item.AppName)
		}
		offset += size
	}
}

func migrate_application_faas() {
	utils.Log.Infof("========== Begin to migrate application_faas ==========")
	var faas []commontypes.ApplicationFaas
	offset = 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&faas).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(faas) == 0 {
			break
		}
		result := destdb.Create(&faas)
		if result.Error != nil {
			utils.Log.Warnf("Batch migrate faas failed - %v", result.Error)
			offset += size
			continue
		}
		utils.Log.Infof("Batch migrate faas success, rows affected=%d", result.RowsAffected)
		offset += size
	}
}

func migrate_application_resource() {
	utils.Log.Infof("========== Begin to migrate application_resource ==========")
	var resource []commontypes.ApplicationResource
	offset = 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&resource).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(resource) == 0 {
			break
		}
		result := destdb.Create(&resource)
		if result.Error != nil {
			utils.Log.Warnf("Batch migrate application_resource failed - %v", result.Error)
			offset += size
			continue
		}
		utils.Log.Infof("Batch migrate application_resource success, rows affected=%d", result.RowsAffected)
		offset += size
	}
}

func migrate_ci_trigger() {
	utils.Log.Infof("========== Begin to migrate ci_trigger ==========")
	var triggers []commontypes.CiTrigger
	offset = 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&triggers).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(triggers) == 0 {
			break
		}
		result := destdb.Create(&triggers)
		if result.Error != nil {
			utils.Log.Warnf("Batch migrate ci_trigger failed - %v", result.Error)
			offset += size
			continue
		}
		utils.Log.Infof("Batch migrate ci_trigger success, rows affected=%d", result.RowsAffected)
		offset += size
	}
}

func migrate_git_repository() {
	utils.Log.Infof("========== Begin to migrate git_repository ==========")
	var repos []commontypes.GitRepository
	offset = 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&repos).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(repos) == 0 {
			break
		}
		for _, item := range repos {
			if err := destdb.Save(&item).Error; err != nil {
				utils.Log.Warnf("id=%d git_http_url=%s - %v", item.Id, item.GitHttpUrl, err)
				continue
			}
			utils.Log.Infof("id=%d git_http_url=%s migrate success", item.Id, item.GitHttpUrl)
		}
		offset += size
	}
}

func migrate_helm_repository() {
	utils.Log.Infof("========== Begin to migrate helm_repository ==========")
	var repos []commontypes.HelmRepositories
	offset = 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&repos).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(repos) == 0 {
			break
		}
		for _, item := range repos {
			if err := destdb.Save(&item).Error; err != nil {
				utils.Log.Warnf("id=%d name=%s - %v", item.Id, item.Name, err)
				continue
			}
			utils.Log.Infof("id=%d name=%s migrate success", item.Id, item.Name)
		}
		offset += size
	}
}

func migrate_image_repository() {
	utils.Log.Infof("========== Begin to migrate image_repository ==========")
	var repos []commontypes.ImageRepository
	offset = 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&repos).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(repos) == 0 {
			break
		}
		for _, item := range repos {
			if err := destdb.Save(&item).Error; err != nil {
				utils.Log.Warnf("id=%d repo_url=%s - %v", item.Id, item.RepoUrl, err)
				continue
			}
			utils.Log.Infof("id=%d repo_url=%s migrate success", item.Id, item.RepoUrl)
		}
		offset += size
	}
}

func migrate_oauth2() {
	utils.Log.Infof("========== Begin to migrate oauth2 ==========")
	var oauth []commontypes.Oauth2
	offset = 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&oauth).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(oauth) == 0 {
			break
		}
		for _, item := range oauth {
			if err := destdb.Save(&item).Error; err != nil {
				utils.Log.Warnf("id=%d name=%s - %v", item.Id, item.Name, err)
				continue
			}
			utils.Log.Infof("id=%d name=%s migrate success", item.Id, item.Name)
		}
		offset += size
	}
}

func migrate_settings() {
	utils.Log.Infof("========== Begin to migrate settings ==========")
	var settings []commontypes.Settings
	offset = 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&settings).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(settings) == 0 {
			break
		}
		result := destdb.Create(&settings)
		if result.Error != nil {
			utils.Log.Warnf("Batch migrate settings failed - %v", result.Error)
			offset += size
			continue
		}
		utils.Log.Infof("Batch migrate settings success, rows affected=%d", result.RowsAffected)
		offset += size
	}
}

func migrate_task_history() {
	utils.Log.Infof("========== Begin to migrate task_history ==========")
	var history []commontypes.TaskHistory
	offset = 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&history).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(history) == 0 {
			break
		}
		result := destdb.Create(&history)
		if result.Error != nil {
			utils.Log.Warnf("Batch migrate task_history failed - %v", result.Error)
			offset += size
			continue
		}
		utils.Log.Infof("Batch migrate task_history success, rows affected=%d", result.RowsAffected)
		offset += size
	}
}

func migrate_task_history_workload() {
	utils.Log.Infof("========== Begin to migrate task_history_workload ==========")
	var workload []commontypes.TaskHistoryWorkload
	offset = 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&workload).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(workload) == 0 {
			break
		}
		result := destdb.Create(&workload)
		if result.Error != nil {
			utils.Log.Warnf("Batch migrate task_history_workload failed - %v", result.Error)
			offset += size
			continue
		}
		utils.Log.Infof("Batch migrate task_history_workload success, rows affected=%d", result.RowsAffected)
		offset += size
	}
}

func migrate_tenant() {
	utils.Log.Infof("========== Begin to migrate tenant ==========")
	var tenant []commontypes.Tenant
	offset = 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&tenant).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(tenant) == 0 {
			break
		}
		for _, item := range tenant {
			if err := destdb.Save(&item).Error; err != nil {
				utils.Log.Warnf("id=%d tenant_name=%s - %v", item.Id, item.TenantName, err)
				continue
			}
			utils.Log.Infof("id=%d tenant_name=%s migrate success", item.Id, item.TenantName)
		}
		offset += size
	}
}

func migrate_token() {
	utils.Log.Infof("========== Begin to migrate tokens ==========")
	var tokens []commontypes.Tokens
	offset = 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&tokens).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(tokens) == 0 {
			break
		}
		for _, item := range tokens {
			if err := destdb.Save(&item).Error; err != nil {
				utils.Log.Warnf("id=%d user_name=%s - %v", item.Id, item.Username, err)
				continue
			}
			utils.Log.Infof("id=%d user_name=%s migrate success", item.Id, item.Username)
		}
		offset += size
	}
}

func migrate_user() {
	utils.Log.Infof("========== Begin to migrate user ==========")
	var users []commontypes.User
	offset = 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&users).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(users) == 0 {
			break
		}
		result := destdb.Create(&users)
		if result.Error != nil {
			utils.Log.Warnf("Batch migrate user failed - %v", result.Error)
			offset += size
			continue
		}
		utils.Log.Infof("Batch migrate user success, rows affected=%d", result.RowsAffected)
		offset += size
	}
}

func migrate_yaml_template() {
	utils.Log.Infof("========== Begin to migrate yaml_template ==========")
	var templates []commontypes.YamlTemplate
	offset = 0
	for {
		if err := srcdb.Limit(size).Offset(offset).Find(&templates).Error; err != nil {
			utils.Log.Warn(err)
			return
		}
		if len(templates) == 0 {
			break
		}
		result := destdb.Create(&templates)
		if result.Error != nil {
			utils.Log.Warnf("Batch migrate yaml_template failed - %v", result.Error)
			offset += size
			continue
		}
		utils.Log.Infof("Batch migrate yaml_template success, rows affected=%d", result.RowsAffected)
		offset += size
	}
}
