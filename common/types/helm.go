package types

import (
	"time"
)

type HelmRepositories struct {
	Id                    int    `json:"id" gorm:"primaryKey,autoIncrement"`
	Name                  string `json:"name" gorm:"size:50;unique"`
	URL                   string `json:"url"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	CertFile              string `json:"cert_file"`
	KeyFile               string `json:"key_file"`
	CAFile                string `json:"ca_file"`
	InsecureSkipTLSverify bool   `json:"insecure_skip_tlsverify"`
	PassCredentialsAll    bool   `json:"pass_credentials_all"`
	Tenant                string `json:"tenant"`
}

type ChartListResponse struct {
	RepoUrl      string `json:",omitempty"`
	ChartName    string
	ChartVersion string
	AppVersion   string
	Icon         string `json:",omitempty"`
	Description  string
}

type ChartInstallRequest struct {
	RepoURL      string
	ChartName    string
	ChartVersion string
	Namespace    string
	ReleaseName  string
	Values       map[string]interface{}
}

type ReleaseElement struct {
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	Revision     string `json:"revision"`
	Updated      string `json:"updated"`
	Status       string `json:"status"`
	Chart        string `json:"chart"`
	ChartName    string `json:"chart_name"`
	ChartVersion string `json:"chart_version"`
	AppVersion   string `json:"app_version"`
}

type ReleaseHistoryInfo struct {
	Revision    int       `json:"revision"`
	Updated     time.Time `json:"updated"`
	Status      string    `json:"status"`
	Chart       string    `json:"chart"`
	AppVersion  string    `json:"app_version"`
	Description string    `json:"description"`
}
