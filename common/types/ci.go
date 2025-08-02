package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type CiTrigger struct {
	Id              int         `json:"id" gorm:"primaryKey,autoIncrement"`
	TriggerName     string      `json:"trigger_name" gorm:"size:300"`
	AppId           int         `json:"app_id" gorm:"uniqueIndex:idx_app_id"`
	Application     Application `json:"application" gorm:"foreignKey:AppId;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	TriggerPath     StringList  `json:"trigger_path" gorm:"type:json"`
	GitHttpUrl      string      `json:"git_http_url" gorm:"size:300"`
	TriggerType     string      `json:"trigger_type" gorm:"size:100"`
	TriggerEndpoint string      `json:"trigger_endpoint" gorm:"size:300"`
	TriggerBody     *string     `json:"trigger_body,omitempty" gorm:"type:json"`
	TriggerEvent    string      `json:"trigger_event" gorm:"size:100"`
	ScanPath        string      `json:"scan_path" gorm:"size:300"`
	RefPattern      string      `json:"ref_pattern"`
	MatchLabels     StringList  `json:"match_labels" gorm:"type:json"`
	PreCheck        PreCheck    `json:"pre_check" gorm:"type:json"`
	UniqInstance    bool        `json:"uniq_instance"`
	Tenant          string      `json:"tenant" gorm:"size:50"`
	UpdateAt        time.Time   `json:"update_at"`
}

type PreCheck struct {
	Enable       bool   `json:"enable"`
	RefPattern   string `json:"ref_pattern"`
	TriggerEvent string `json:"trigger_event"`
}

func (p PreCheck) Value() (driver.Value, error) {
	b, err := json.Marshal(p)
	return string(b), err
}

func (p *PreCheck) Scan(value interface{}) error {
	switch v := value.(type) {
	case []uint8:
		return json.Unmarshal(v, &p)
	default:
		return json.Unmarshal([]byte(value.(string)), &p)
	}
}

type GitRepository struct {
	Id             int    `json:"id" gorm:"primaryKey,autoIncrement"`
	GitHttpUrl     string `json:"git_http_url" gorm:"size:300;uniqueIndex:idx_git_repo"`
	AccessToken    string `json:"access_token" gorm:"size:100;uniqueIndex:idx_git_repo"`
	WebhookBaseURL string `json:"webhook_base_url" gorm:"size:300"`
	Tenant         string `json:"tenant" gorm:"size:50;uniqueIndex:idx_git_repo"`
}

type LogGitWebhook struct {
	Id         int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Timestamp  time.Time `json:"timestamp"`
	ObjectKind string    `json:"object_kind" gorm:"size:20"`
	GitHttpUrl string    `json:"git_http_url" gorm:"size:300"`
	MergeIID   int64     `json:"merge_iid"`
	Revision   string    `json:"revision" gorm:"size:100"`
	Committer  string    `json:"committer" gorm:"size:100"`
	Action     string    `json:"action" gorm:"size:50"`
	CommitId   string    `json:"commit_id" gorm:"size:100"`
	Ignored    bool      `json:"ignored"`
	TriggerApp string    `json:"trigger_app"`
	Reason     string    `json:"reason"`
	Body       string    `json:"body" gorm:"type:json"`
}
