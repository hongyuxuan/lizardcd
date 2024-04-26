package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

type GetDataReq struct {
	Tablename string `path:"tablename"`
	Page      int    `form:"page,default=1"`
	Size      int    `form:"size,default=20"`
	Search    string `form:"search,optional"`
	Filter    string `form:"filter,optional"`
	Range     string `form:"range,optional"`
	Sort      string `form:"sort,optional"`
}

type ApplicationTemplate struct {
	Id        int       `yaml:"id" gorm:"primaryKey,autoIncrement"`
	Name      string    `yaml:"name" gorm:"unique"`
	Content   string    `yaml:"content"`
	Variables string    `yaml:"variables"`
	UpdateAt  time.Time `yaml:"update_at"`
}

type User struct {
	Id       int    `gorm:"primaryKey,autoIncrement"`
	Username string `gorm:"size:50;unique"`
	Password string `gorm:"size:100"`
	UpdateAt time.Time
}

type ImageRepository struct {
	Id           int    `json:"id" gorm:"primaryKey,autoIncrement"`
	RepoUrl      string `json:"repo_url" gorm:"size:300"`
	RepoAccount  string `json:"repo_account" gorm:"size:50"`
	RepoPassword string `json:"repo_password" gorm:"size:100"`
	RepoType     string `json:"repo_type" gorm:"size:20"`
}

func (r ImageRepository) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *ImageRepository) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &r)
}

type Application struct {
	Id        int             `json:"id" gorm:"primaryKey,autoIncrement"`
	AppName   string          `json:"app_name" gorm:"size:300;unique"`
	Repo      ImageRepository `json:"repo" gorm:"type:json"`
	RepoName  string          `json:"repo_name" gorm:"size:50"`
	ImageName string          `json:"image_name" gorm:"size:300"`
	Workload  WorkLoadList    `json:"workload" gorm:"type:json"`
	UpdateAt  time.Time       `json:"update_at"`
}

type WorkLoad struct {
	Cluster       string `json:"cluster"`
	Namespace     string `json:"namespace"`
	WorkloadName  string `json:"workload_name"`
	ContainerName string `json:"container_name"`
	WorkloadType  string `json:"workload_type,optional"`
}

type WorkLoadList []WorkLoad

func (w WorkLoadList) Value() (driver.Value, error) {
	return json.Marshal(w)
}

func (w *WorkLoadList) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &w)
}

type LoginRes struct {
	AccessToken  string `json:"access_token"`
	AccessExpire int64  `json:"access_expire"`
	RefreshAfter int64  `json:"refresh_after"`
}

type JfrogFileItem struct {
	Uri          string `json:"uri"`
	Size         int64  `json:"size"`
	LastModified string `json:"lastModified"`
	Folder       bool   `json:"folder"`
}

type HarborFileItem struct {
	Type string `json:"type"`
	Size int64  `json:"size"`
	Tags []struct {
		Name     string `json:"name"`
		PushTime string `json:"push_time"`
	} `json:"tags"`
}

type ArtifactListRes struct {
	ArtifactUrl  string `json:"artifact_url"`
	LastModified string `json:"last_modified"`
	Tag          string `json:"tag"`
}
