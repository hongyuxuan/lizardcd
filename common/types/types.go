package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type TraceIDKey struct{}

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

type StringList []string

func (s StringList) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	str := string(b)
	if str == "null" {
		str = "[]"
	}
	return str, err
}

func (s *StringList) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &s)
}

type StringMap map[string]string

func (s StringMap) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	str := string(b)
	if str == "null" {
		str = "{}"
	}
	return str, err
}

func (s *StringMap) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &s)
}

type GetDataReq struct {
	Tablename string `path:"tablename"`
	Page      int    `form:"page,default=1"`
	Size      int    `form:"size,default=20"`
	Search    string `form:"search,optional"`
	Filter    string `form:"filter,optional"`
	Range     string `form:"range,optional"`
	Preload   bool   `form:"preload,optional"`
	Sort      string `form:"sort,optional"`
}

type ApplicationTemplate struct {
	Id        int       `yaml:"id" gorm:"primaryKey,autoIncrement"`
	Name      string    `yaml:"name" gorm:"unique"`
	Content   string    `yaml:"content"`
	Variables string    `yaml:"variables"`
	Tenant    string    `yaml:"tenant"`
	UpdateAt  time.Time `yaml:"update_at"`
}

type User struct {
	Id       int       `json:"id" gorm:"primaryKey,autoIncrement"`
	Username string    `json:"username" gorm:"size:50;unique"`
	Password string    `json:",omitempty" gorm:"size:100"`
	Role     string    `json:"role" gorm:"size:50"`
	Tenant   string    `json:"tenant" gorm:"size:50"`
	UpdateAt time.Time `json:"update_at"`
}

type Tenant struct {
	Id         int       `json:"id" gorm:"primaryKey,autoIncrement"`
	TenantName string    `json:"tenant_name" gorm:"size:50;unique"`
	Namespaces string    `json:"namespaces"`
	UpdateAt   time.Time `json:"update_at"`
}

type ImageRepository struct {
	Id           int    `json:"id" gorm:"primaryKey,autoIncrement"`
	RepoUrl      string `json:"repo_url" gorm:"size:300;uniqueIndex:idx_repo"`
	RepoAccount  string `json:"repo_account" gorm:"size:50;uniqueIndex:idx_repo"`
	RepoPassword string `json:"repo_password" gorm:"size:100"`
	RepoType     string `json:"repo_type" gorm:"size:20"`
	Tenant       string `json:"tenant" gorm:"size:50;uniqueIndex:idx_repo"`
}

func (r ImageRepository) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *ImageRepository) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &r)
}

type Application struct {
	Id                   int             `json:"id" gorm:"primaryKey,autoIncrement"`
	AppName              string          `json:"app_name" gorm:"size:300;unique"`
	DeployType           string          `json:"deploy_type" gorm:"size:10"`
	Repo                 ImageRepository `json:"repo" gorm:"type:json"`
	RepoName             string          `json:"repo_name" gorm:"size:50"`
	ImageName            string          `json:"image_name" gorm:"size:300"`
	Workload             WorkLoadList    `json:"workload" gorm:"type:json"`
	EnableTrafficControl bool            `json:"enable_traffic_control"`
	TrafficPolicy        string          `json:"traffic_policy"`
	Tenant               string          `json:"tenant" gorm:"size:50"`
	Tags                 StringList      `json:"tags" gorm:"type:json"`
	ExtraVars            string          `json:"extra_vars" gorm:"type:json"`
	UpdateAt             time.Time       `json:"update_at"`
}

func (a Application) Tablename() string {
	return "application"
}

type ApplicationList []Application

type WorkLoad struct {
	Cluster       string        `json:"cluster,omitempty"`
	Namespace     string        `json:"namespace,omitempty"`
	WorkloadName  string        `json:"workload_name"`
	ContainerName string        `json:"container_name,omitempty"`
	WorkloadType  string        `json:"workload_type,optional"`
	Version       string        `json:"version,optional,omitempty"`
	Weight        int           `json:"weight,optional,omitempty"`
	Headers       []MatchHeader `json:"headers,optional,omitempty"`
	ArtifactUrl   string        `json:"artifact_url,optional,omitempty"`
	Enable        bool          `json:"enable"`
}

type MatchHeader struct {
	Key       string `json:"key"`
	MatchType string `json:"match_type"`
	Value     string `json:"value"`
}

func (w WorkLoad) Value() (driver.Value, error) {
	b, err := json.Marshal(w)
	return string(b), err
}

func (w *WorkLoad) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &w)
}

type WorkLoadList []WorkLoad

func (w WorkLoadList) Value() (driver.Value, error) {
	return json.Marshal(w)
}

func (w *WorkLoadList) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &w)
}

type Settings struct {
	Id           int    `json:"id" gorm:"primaryKey;autoIncrement"`
	SettingKey   string `json:"setting_key" gorm:"uniqueIndex:idx_key"`
	SettingValue string `json:"setting_value"`
	Tenant       string `json:"tenant" gorm:"size:50;uniqueIndex:idx_key"`
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

type DockerHubImageItem struct {
	Name        string `json:"name"`
	FullSize    int64  `json:"full_size"`
	LastUpdated string `json:"last_updated"`
}

type ArtifactListRes struct {
	ArtifactUrl  string `json:"artifact_url"`
	LastModified string `json:"last_modified"`
	Tag          string `json:"tag"`
}

type TaskHistory struct {
	Id                   string                `json:"id" gorm:"size:100;primaryKey"`
	AppName              string                `json:"app_name" gorm:"size:50"`
	TaskType             string                `json:"task_type" gorm:"size:50"`
	Success              sql.NullBool          `json:"success"`
	ErrMessage           string                `json:"err_message"`
	Status               string                `json:"status" gorm:"size:50"`
	Tenant               string                `json:"tenant" gorm:"size:50"`
	TriggerType          string                `json:"trigger_type" gorm:"size:50"`
	InitAt               sql.NullTime          `json:"init_at"`
	StartAt              sql.NullTime          `json:"start_at"`
	FinishAt             sql.NullTime          `json:"finish_at"`
	Expire               string                `json:"expire" gorm:"size:50"`
	Labels               StringList            `json:"labels" gorm:"type:json"`
	TaskHistoryWorkloads []TaskHistoryWorkload `json:"workloads" gorm:"foreignKey:TaskHistoryId;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

type TaskHistoryWorkload struct {
	Id            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Workload      WorkLoad  `json:"workload" gorm:"type:json"`
	Status        string    `json:"status" gorm:"type:json"`
	ErrMessage    string    `json:"err_message"`
	TaskHistoryId string    `json:"task_history_id" gorm:"size:100"`
	UpdateAt      time.Time `json:"update_at"`
}
