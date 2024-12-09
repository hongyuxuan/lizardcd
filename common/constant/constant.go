package constant

const REPO_TYPE_ARTIFACTORY = "Artifactory"
const REPO_TYPE_HARBOR = "Harbor"
const REPO_TYPE_DOCKERHUB = "DockerHub"
const REPO_TYPE_S3 = "S3"

const K8S_RESOURCE_TYPE_DEPLOYMENTS = "deployments"
const K8S_RESOURCE_TYPE_STATEFULSETS = "statefulsets"

const ISTIO_CRD_TYPE_DESTINATIONRULE = "destinationrules"
const ISTIO_CRD_TYPE_VIRTUALSERVICE = "virtualservices"
const ISTIO_CRD_TYPE_GATEWAY = "gateways"

const ROLE_ADMIN = "admin"
const ROLE_READWRITE = "readwrite"
const ROLE_READONLY = "readonly"

const TASK_TYPE_ROLLOUT = "rollout"
const TASK_TYPE_DEPLOY = "deploy"
const TASK_TYPE_SYNCHRONIZE = "synchronize"
const TASK_STATUS_INITIALIZE = "initialize"
const TASK_STATUS_WAITING = "waiting"
const TASK_STATUS_TERMINATED = "terminated"
const TASK_STATUS_RUNNING = "running"
const TASK_STATUS_FINISHED = "finished"
const TASK_TRIGGER_TYPE_CRON = "定时调度"

const WS_MSG_TYPE_PODLOG = "LIZARDCD_WS_PODLOG"

const APP_TRAFFIC_POLICY_WEIGHT = "weight"
const APP_TRAFFIC_POLICY_HEADER = "header"
const APP_SYNC_TYPE_AUTO = "自动"
const APP_SYNC_TYPE_MANUAL = "手动"
