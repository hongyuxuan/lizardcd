package types

import (
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

/*********** map[string]string ***********/
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

/*********** map[string]interface{} ***********/
type InterfaceMap map[string]interface{}

func (s InterfaceMap) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	str := string(b)
	if str == "null" {
		str = "{}"
	}
	return str, err
}

func (s *InterfaceMap) Scan(value interface{}) error {
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

type MatchHeader struct {
	Key       string `json:"key"`
	MatchType string `json:"match_type"`
	Value     string `json:"value"`
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

type DolphinschedulerSetting struct {
	BaseUrl          string `json:"base_url"`
	Token            string `json:"token"`
	Project          string `json:"project"`
	LizardcdJwtToken string `json:"lizardcd_jwt_token"`
}

type ListResult struct {
	Total   int         `json:"total"`
	Results interface{} `json:"results"`
}

type DolphinschedulerResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	Failed  bool        `json:"failed"`
	Success bool        `json:"success"`
}
