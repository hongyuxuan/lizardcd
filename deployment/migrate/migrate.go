package main

import (
	"flag"
	"io"
	"os"
	"time"

	"github.com/hongyuxuan/lizardcd/common/types"
	"github.com/hongyuxuan/lizardcd/common/utils"
	"gopkg.in/yaml.v2"
)

var dbfile = flag.String("d", "./lizardcd.db", "sqlite database file")
var level = flag.String("l", "info", "sqlite database file")

func main() {
	flag.Parse()

	db := utils.NewSQLite(*dbfile, *level)

	var file *os.File
	var err error

	// create table `yaml_template``
	db.AutoMigrate(&types.YamlTemplate{})
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
		if err = db.Save(&yamlTemplate).Error; err != nil {
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
		if err = db.Save(&yamlTemplate).Error; err != nil {
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
		if err = db.Save(&yamlTemplate).Error; err != nil {
			utils.Log.Warn(err)
			continue
		}
		utils.Log.Infof("saved name=\"%s\" into yaml_template success", yamlTemplate.Name)
	}

	// create table `user`
	db.AutoMigrate(&types.User{})
	generatedPassword := utils.GenerateRandomString(10)
	if err = utils.AddUser("admin", generatedPassword, "admin", "admin", db); err != nil {
		utils.Log.Warnf("failed to migrate table user: %v", err)
		db.Model(&types.User{}).Where("username = ?", "admin").Update("tenant", "admin")
		db.Model(&types.User{}).Where("username = ?", "admin").Update("role", "admin")
	} else {
		utils.Log.Infof("password of user admin is: %s , please modified it when you first login", generatedPassword)
	}

	// create table `tenant`
	db.AutoMigrate(&types.Tenant{})
	if err = db.Save(&types.Tenant{
		TenantName: "admin",
		Namespaces: "[]",
	}).Error; err != nil {
		utils.Log.Warn(err)
	}

	// create table `image_repository`
	if err = db.AutoMigrate(&types.ImageRepository{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `application`
	if err = db.AutoMigrate(&types.Application{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `settings`
	var tenants []types.Tenant
	db.Find(&tenants)
	// db.Migrator().DropTable(&types.Settings{})
	db.AutoMigrate(&types.Settings{})
	db.Save(&types.Settings{
		SettingKey:   "default_tekton",
		SettingValue: "{}",
		Tenant:       "admin",
	})
	db.Save(&types.Settings{
		SettingKey:   "server_url",
		SettingValue: "",
		Tenant:       "admin",
	})
	for _, t := range tenants {
		utils.AddSettings(t.TenantName, db)
	}

	// create table `helm_repositories`
	var repos []types.HelmRepositories
	if db.Migrator().HasTable(&types.HelmRepositories{}) {
		db.Find(&repos)
	}
	db.Migrator().DropTable(&types.HelmRepositories{})
	db.AutoMigrate(&types.HelmRepositories{})
	for i, r := range repos {
		r.Id = i
		if err := db.Create(&r).Error; err != nil {
			utils.Log.Warn(err)
		}
	}

	// create table `task_history`
	if err = db.AutoMigrate(&types.TaskHistory{}); err != nil {
		utils.Log.Warn(err)
	}
	if err = db.AutoMigrate(&types.TaskHistoryWorkload{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `oauth2`
	if err = db.AutoMigrate(&types.Oauth2{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `application_faas`
	if err = db.AutoMigrate(&types.ApplicationFaas{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `application_resource`
	if err = db.AutoMigrate(&types.ApplicationResource{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `ci_trigger`
	if err = db.AutoMigrate(&types.CiTrigger{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `git_repository`
	if err = db.AutoMigrate(&types.GitRepository{}); err != nil {
		utils.Log.Warn(err)
	}

	// create table `tokens`
	if err = db.AutoMigrate(&types.Tokens{}); err != nil {
		utils.Log.Warn(err)
	}
}
