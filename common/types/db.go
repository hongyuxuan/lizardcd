package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type YamlTemplate struct {
	Id        int       `yaml:"id" json:"id" gorm:"primaryKey,autoIncrement"`
	Name      string    `yaml:"name" json:"name" gorm:"unique"`
	Type      string    `yaml:"type" json:"type" gorm:"size:50"`
	Content   string    `yaml:"content" json:"content"`
	Variables string    `yaml:"variables" json:"variables"`
	Tenant    string    `yaml:"tenant" json:"tenant" gorm:"size:50"`
	UpdateAt  time.Time `json:"update_at"`
}

type User struct {
	Id       int       `json:"id" gorm:"primaryKey,autoIncrement"`
	Username string    `json:"username" gorm:"size:50;unique"`
	Password string    `json:",omitempty" gorm:"size:100"`
	Role     string    `json:"role" gorm:"size:50"`
	Tenant   string    `json:"tenant" gorm:"size:50"`
	Profile  StringMap `json:"profile" gorm:"type:json"`
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
	Id                   int                   `json:"id" gorm:"primaryKey,autoIncrement"`
	AppName              string                `json:"app_name" gorm:"size:300;unique"`
	DeployType           string                `json:"deploy_type" gorm:"size:10"`
	GitHttpUrl           string                `json:"git_http_url" gorm:"size:300"`
	RepoId               int                   `json:"repo_id"`
	RepoName             string                `json:"repo_name" gorm:"size:50"`
	ImageName            string                `json:"image_name" gorm:"size:300"`
	Workload             WorkloadList          `json:"workload" gorm:"type:json"`
	EnableTrafficControl bool                  `json:"enable_traffic_control"`
	TrafficPolicy        string                `json:"traffic_policy"`
	Tenant               string                `json:"tenant" gorm:"size:50"`
	Tags                 StringList            `json:"tags" gorm:"type:json"`
	ExtraVars            string                `json:"extra_vars" gorm:"type:json"`
	TemplateId           int                   `json:"template_id,omitempty"`
	EnableBuild          bool                  `json:"enable_build"`
	BuildScript          string                `json:"build_script,omitempty"`
	VersionScript        string                `json:"version_script,omitempty"`
	GitOps               *GitOps               `json:"gitops,omitempty" gorm:"type:json"`
	UpdateAt             time.Time             `json:"update_at"`
	Template             *YamlTemplate         `json:"template,omitempty" gorm:"foreignKey:TemplateId;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Faas                 *ApplicationFaas      `json:"faas,omitempty" gorm:"foreignKey:ApplicationId;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Resource             []ApplicationResource `json:"resource,omitempty" gorm:"foreignKey:ApplicationId;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

type ApplicationList []Application

type Workload struct {
	Cluster       string        `json:"cluster,omitempty"`
	Namespace     string        `json:"namespace,omitempty"`
	WorkloadName  string        `json:"workload_name"`
	ContainerName string        `json:"container_name,omitempty"`
	WorkloadType  string        `json:"workload_type,optional,omitempty"`
	Version       string        `json:"version,optional,omitempty"`
	Weight        int           `json:"weight,optional,omitempty"`
	Headers       []MatchHeader `json:"headers,optional,omitempty"`
	ArtifactUrl   string        `json:"artifact_url,optional,omitempty"`
	// Revision      string        `json:"revision,optional,omitempty"`
	Enable bool `json:"enable"`
}

func (w Workload) Value() (driver.Value, error) {
	b, err := json.Marshal(w)
	return string(b), err
}

func (w *Workload) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &w)
}

type WorkloadList []Workload

func (w WorkloadList) Value() (driver.Value, error) {
	b, err := json.Marshal(w)
	return string(b), err
}
func (w *WorkloadList) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &w)
}

type GitOps struct {
	GitRevision   string `json:"git_revision"`
	CommitId      string `json:"commit_id"`
	Path          string `json:"path"`
	Include       string `json:"include,omitempty"`
	Exclude       string `json:"exclude,omitempty"`
	SyncType      string `json:"sync_type"`
	Cron          string `json:"cron,omitempty"`
	HostPortId    string `json:"host_port_id,omitempty"`
	PruneOnDelete bool   `json:"prune_on_delete,omitempty"`
	UseDS         bool   `json:"use_ds,omitempty"`
	ProcessCode   int64  `json:"process_code,omitempty"`
	CronId        int64  `json:"cron_id,omitempty"`
}

func (g GitOps) Value() (driver.Value, error) {
	b, err := json.Marshal(g)
	return string(b), err
}

func (g *GitOps) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &g)
}

type ApplicationResource struct {
	ApplicationId   int    `json:"application_id" gorm:"uniqueIndex:idx_app_resource"`
	Cluster         string `json:"cluster" gorm:"size:50;uniqueIndex:idx_app_resource"`
	Namespace       string `json:"namespace" gorm:"size:50;uniqueIndex:idx_app_resource"`
	ResourceType    string `json:"resource_type" gorm:"size:50;uniqueIndex:idx_app_resource"`
	ResourceName    string `json:"resource_name" gorm:"size:50;uniqueIndex:idx_app_resource"`
	DesiredManifest string `json:"desired_manifest"`
	LastManifest    string `json:"last_manifest"`
	LiveManifest    string `json:"live_manifest,omitempty"`
}

type ApplicationFaas struct {
	Id            int             `json:"id" gorm:"primaryKey,autoIncrement"`
	Cluster       string          `json:"cluster" gorm:"size:50"`
	Namespace     string          `json:"namespace" gorm:"size:50"`
	Template      string          `json:"content"`
	Variables     InterfaceMap    `json:"variables,omitempty" gorm:"type:json"`
	Versions      FaasVersionList `json:"versions" gorm:"type:json"`
	ApplicationId int             `json:"application_id"`
}

type FaasVersion struct {
	CpuUsage  int    `json:"cpu_usage,omitempty"`
	MemUsage  int    `json:"mem_usage,omitempty"`
	EnableHpa bool   `json:"enable_hpa,omitempty"`
	MaxPod    int    `json:"max_pod,omitempty"`
	MinPod    int    `json:"min_pod,omitempty"`
	Version   string `json:"version,omitempty"`
	Weight    int    `json:"weight,omitempty"`
}

type FaasVersionList []FaasVersion

func (f FaasVersionList) Value() (driver.Value, error) {
	b, err := json.Marshal(f)
	return string(b), err
}

func (f *FaasVersionList) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), &f)
}

type Settings struct {
	Id           int    `json:"id" gorm:"primaryKey;autoIncrement"`
	SettingKey   string `json:"setting_key" gorm:"uniqueIndex:idx_key"`
	SettingValue string `json:"setting_value"`
	Tenant       string `json:"tenant" gorm:"size:50;uniqueIndex:idx_key"`
}

type TaskHistory struct {
	Id                   string                `json:"id" gorm:"size:100;primaryKey"`
	AppName              string                `json:"app_name" gorm:"index:idx_history_app;size:50"`
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
	Workload      Workload  `json:"workload" gorm:"type:json"`
	Status        string    `json:"status" gorm:"type:json"`
	ErrMessage    string    `json:"err_message"`
	TaskHistoryId string    `json:"task_history_id" gorm:"size:100"`
	UpdateAt      time.Time `json:"update_at"`
}

type Oauth2 struct {
	Id             int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string `json:"name" gorm:"size:50"`
	ClientId       string `json:"client_id" gorm:"size:50"`
	ClientSecret   string `json:"client_secret" gorm:"size:50"`
	AuthorizeUrl   string `json:"authorized_url"`
	RedirectUrl    string `json:"redirect_url"`
	TokenUrl       string `json:"token_url"`
	UserinfoUrl    string `json:"userinfo_url"`
	CallbackUrl    string `json:"callback_url"`
	UserJsonpath   string `json:"user_jsonpath"`
	AvatarJsonpath string `json:"avatar_jsonpath"`
	HttpProxy      string `json:"http_proxy"`
}

type Tokens struct {
	Id       int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId   int       `json:"user_id"`
	Username string    `json:"username" gorm:"size:50"`
	User     User      `json:"user,omitempty" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	JwtToken string    `json:"jwt_token"`
	CreateAt time.Time `json:"create_at"`
	ExpireAt time.Time `json:"expire_at"`
}
