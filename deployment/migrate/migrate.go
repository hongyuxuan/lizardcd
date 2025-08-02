package main

import (
	"io"
	"os"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

var (
	dbtype   = kingpin.Flag("db-type", "database type(sqlite|tidb|mysql|mariadb)").Short('t').Default("sqlite").String()
	dbfile   = kingpin.Flag("db-file", "database file, only for sqlite").Short('d').Default("./lizardcd.db").String()
	dbHost   = kingpin.Flag("db-host", "database host").Short('H').Default("localhost").String()
	dbPort   = kingpin.Flag("db-port", "database port").Short('P').Default("3306").Int()
	username = kingpin.Flag("db-username", "database username").Short('u').Default("lizardcd").String()
	password = kingpin.Flag("db-password", "database password").Short('p').String()
	db       = kingpin.Flag("db", "database").Default("lizardcd").String()
	logLevel = kingpin.Flag("log.level", "Log level.").Default("info").String()
)

func main() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	var database *gorm.DB
	if *dbtype == "tidb" {
		database = utils.NewTidb(*dbHost, *username, *password, *db, *logLevel, *dbPort)
	} else {
		database = utils.NewSQLite(*dbfile, *logLevel)
	}
	utils.InitLogger(*logLevel)

	var file *os.File
	var err error

	// create table `yaml_template`
	database.AutoMigrate(&types.YamlTemplate{})
	// application_template
	if file, err = os.Open("manifests/application_template.yaml"); err != nil {
		utils.Log.Fatal(err)
	}
	dec := yaml.NewDecoder(file)
	for {
		var yamlTemplate types.YamlTemplate
		err = dec.Decode(&yamlTemplate)
		if err == io.EOF {
			break
		}
		if err != nil {
			utils.Log.Warn(err)
		}
		yamlTemplate.UpdateAt = time.Now()
		yamlTemplate.Type = "application"
		yamlTemplate.Tenant = "admin"
		if err = database.Save(&yamlTemplate).Error; err != nil {
			utils.Log.Warn(err)
			continue
		}
		utils.Log.Infof("saved name=\"%s\" into yaml_template success", yamlTemplate.Name)
	}

	// tekton_template
	if file, err = os.Open("manifests/tekton_template.yaml"); err != nil {
		utils.Log.Fatal(err)
	}
	dec = yaml.NewDecoder(file)
	for {
		var yamlTemplate types.YamlTemplate
		err = dec.Decode(&yamlTemplate)
		if err == io.EOF {
			break
		}
		if err != nil {
			utils.Log.Warn(err)
		}
		yamlTemplate.UpdateAt = time.Now()
		yamlTemplate.Tenant = "admin"
		if err = database.Save(&yamlTemplate).Error; err != nil {
			utils.Log.Warn(err)
			continue
		}
		utils.Log.Infof("saved name=\"%s\" into yaml_template success", yamlTemplate.Name)
	}

	// processdefine_template
	if file, err = os.Open("manifests/processdefine.yaml"); err != nil {
		utils.Log.Fatal(err)
	}
	dec = yaml.NewDecoder(file)
	for {
		var yamlTemplate types.YamlTemplate
		err = dec.Decode(&yamlTemplate)
		if err == io.EOF {
			break
		}
		if err != nil {
			utils.Log.Warn(err)
		}
		yamlTemplate.UpdateAt = time.Now()
		yamlTemplate.Tenant = "admin"
		if err = database.Save(&yamlTemplate).Error; err != nil {
			utils.Log.Warn(err)
			continue
		}
		utils.Log.Infof("saved name=\"%s\" into yaml_template success", yamlTemplate.Name)
	}

	// create table `user`
	database.AutoMigrate(&types.User{})
	generatedPassword := utils.GenerateRandomString(10)
	if err = utils.AddUser("admin", generatedPassword, "admin", "admin", database); err != nil {
		utils.Log.Warnf("failed to migrate table user: %v", err)
		database.Model(&types.User{}).Where("username = ?", "admin").Update("tenant", "admin")
		database.Model(&types.User{}).Where("username = ?", "admin").Update("role", "admin")
	} else {
		utils.Log.Infof("password of user admin is: %s , please modified it when you first login", generatedPassword)
	}

	// create table `tenant`
	database.AutoMigrate(&types.Tenant{})
	if err = database.Save(&types.Tenant{
		TenantName: "admin",
		Namespaces: "[]",
		UpdateAt:   time.Now(),
	}).Error; err != nil {
		utils.Log.Warn(err)
	}

	// create table `image_repository`
	if err = database.AutoMigrate(&types.ImageRepository{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `application`
	if err = database.AutoMigrate(&types.Application{}); err != nil {
		utils.Log.Warn(err)
	}
	if database.Migrator().HasColumn(&types.Application{}, "enable_build") {
		database.Migrator().DropColumn(&types.Application{}, "enable_build")
		database.Migrator().DropColumn(&types.Application{}, "build_script")
		database.Migrator().DropColumn(&types.Application{}, "version_script")
		if err = database.Transaction(func(tx *gorm.DB) error {
			if tx.Migrator().HasTable(&types.ApplicationNew{}) {
				tx.Migrator().DropTable(&types.ApplicationNew{})
			}
			tx.AutoMigrate(&types.ApplicationNew{})
			if err := tx.Exec(`
					INSERT INTO application_new (id,app_name,repo_name,image_name,workload,enable_traffic_control,traffic_policy,update_at,tenant,tags,deploy_type,extra_vars,repo_id,git_http_url,git_ops,timeout,auto_build) SELECT id,app_name,repo_name,image_name,workload,enable_traffic_control,traffic_policy,update_at,tenant,tags,deploy_type,extra_vars,repo_id,git_http_url,git_ops,timeout,auto_build FROM application
			`).Error; err != nil {
				return err
			}
			if err := tx.Migrator().DropTable("application"); err != nil {
				return err
			}
			if err := tx.Exec("ALTER TABLE application_new RENAME TO application").Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			utils.Log.Warn(err)
		}
	} else {
		utils.Log.Info("do not need to migrate ApplicationNew")
	}
	if err = database.Model(&types.Application{}).Where("timeout is NULL").Update("timeout", 300).Error; err != nil {
		utils.Log.Warn(err)
	}

	// create table `settings`
	var tenants []types.Tenant
	database.Find(&tenants)
	// database.Migrator().DropTable(&types.Settings{})
	if err = database.AutoMigrate(&types.Settings{}); err != nil {
		utils.Log.Warn(err)
	}
	database.Save(&types.Settings{
		SettingKey:   "default_tekton",
		SettingValue: "{}",
		Tenant:       "admin",
	})
	database.Save(&types.Settings{
		SettingKey:   "server_url",
		SettingValue: "",
		Tenant:       "admin",
	})
	database.Save(&types.Settings{
		SettingKey:   "tekton_source",
		SettingValue: "crd",
		Tenant:       "admin",
	})
	for _, t := range tenants {
		utils.AddSettings(t.TenantName, database)
	}

	// create table `helm_repositories`
	var repos []types.HelmRepositories
	if database.Migrator().HasTable(&types.HelmRepositories{}) {
		database.Find(&repos)
	}
	database.Migrator().DropTable(&types.HelmRepositories{})
	database.AutoMigrate(&types.HelmRepositories{})
	for i, r := range repos {
		r.Id = i
		if err := database.Create(&r).Error; err != nil {
			utils.Log.Warn(err)
		}
	}

	// create table `task_history`
	if err = database.AutoMigrate(&types.TaskHistory{}); err != nil {
		utils.Log.Warn(err)
	}
	if err = database.AutoMigrate(&types.TaskHistoryWorkload{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `oauth2`
	if err = database.AutoMigrate(&types.Oauth2{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `application_faas`
	if err = database.AutoMigrate(&types.ApplicationFaas{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `application_resource`
	if err = database.AutoMigrate(&types.ApplicationResource{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `ci_trigger`
	if err = database.AutoMigrate(&types.CiTrigger{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `git_repository`
	if err = database.AutoMigrate(&types.GitRepository{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `tokens`
	if err = database.AutoMigrate(&types.Tokens{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `agent`
	if err = database.AutoMigrate(&types.Agent{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `log_git_webhook`
	if err = database.AutoMigrate(&types.LogGitWebhook{}); err != nil {
		utils.Log.Warn(err)
	}
}
