package types

import "time"

type CiTrigger struct {
	Id              int         `json:"id" gorm:"primaryKey,autoIncrement"`
	TriggerName     string      `json:"trigger_name" gorm:"size:300"`
	AppId           int         `json:"app_id" gorm:"uniqueIndex:idx_app_id"`
	Application     Application `json:"application" gorm:"foreignKey:AppId;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	TriggerPath     StringList  `json:"trigger_path" gorm:"type:json"`
	GitHttpUrl      string      `json:"git_http_url" gorm:"size:300"`
	TriggerType     string      `json:"trigger_type" gorm:"size:100"`
	TriggerEndpoint string      `json:"trigger_endpoint" gorm:"size:300"`
	TriggerBody     string      `json:"trigger_body" gorm:"type:json"`
	Secret          string      `json:"secret" gorm:"size:300"`
	Tenant          string      `json:"tenant" gorm:"size:50"`
	UpdateAt        time.Time   `json:"update_at"`
}

type GitRepository struct {
	Id             int    `json:"id" gorm:"primaryKey,autoIncrement"`
	GitHttpUrl     string `json:"git_http_url" gorm:"size:300;uniqueIndex:idx_git_repo"`
	AccessToken    string `json:"access_token" gorm:"size:100;uniqueIndex:idx_git_repo"`
	WebhookBaseURL string `json:"webhook_base_url" gorm:"size:300"`
	Secret         string `json:"secret" gorm:"size:300"`
	Tenant         string `json:"tenant" gorm:"size:50;uniqueIndex:idx_git_repo"`
}
