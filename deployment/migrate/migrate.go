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

	// create table `application_template``
	db.AutoMigrate(&types.ApplicationTemplate{})
	if file, err = os.Open("manifests/application_template.yaml"); err != nil {
		utils.Log.Fatal(err)
	}
	dec := yaml.NewDecoder(file)
	for {
		var applicationTemplate types.ApplicationTemplate
		err = dec.Decode(&applicationTemplate)
		if err == io.EOF {
			break
		}
		if err != nil {
			utils.Log.Warn(err)
		}
		applicationTemplate.UpdateAt = time.Now()
		if err = db.Save(&applicationTemplate).Error; err != nil {
			utils.Log.Warn(err)
			continue
		}
		utils.Log.Infof("saved name=\"%s\" into application_template success", applicationTemplate.Name)
	}

	// create table `user`
	db.AutoMigrate(&types.User{})
	generatedPassword := utils.GenerateRandomString(10)
	if err = utils.AddUser("admin", generatedPassword, db); err != nil {
		utils.Log.Warnf("failed to migrate table user: %v", err)
	} else {
		utils.Log.Infof("password of user admin is: %s , please modified it when you first login", generatedPassword)
	}
}
