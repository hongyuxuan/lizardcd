package types

import (
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type LizardAgentRes struct {
	Code int `json:"code"`
	Data []struct {
		ServiceName string `json:"service_name"`
	} `json:"data"`
}

type DeploymentRes struct {
	Code int             `json:"code"`
	Data []v1.Deployment `json:"data"`
}

type StatefulsetRes struct {
	Code int              `json:"code"`
	Data []v1.StatefulSet `json:"data"`
}

type PodRes struct {
	Code int
	Data []corev1.Pod `json:"data"`
}

type EventRes struct {
	Code int
	Data []corev1.Event `json:"data"`
}

type ScaleReq struct {
	Workloads []Workloads `json:"workloads"`
}

type Workloads struct {
	Name     string `json:"name"`
	Replicas int32  `json:"replicas"`
}

type ApplicationRes struct {
	Code int `json:"code"`
	Data struct {
		Total   int                       `json:"total"`
		Results []commontypes.Application `json:"results"`
	} `json:"data"`
}

type TaskHistoriesRes struct {
	Code int `json:"code"`
	Data struct {
		Total   int                       `json:"total"`
		Results []commontypes.TaskHistory `json:"results"`
	} `json:"data"`
}

type TaskHistoryRes struct {
	Code int                     `json:"code"`
	Data commontypes.TaskHistory `json:"data"`
}

type HelmRepoRes struct {
	Code int `json:"code"`
	Data struct {
		Total   int                            `json:"total"`
		Results []commontypes.HelmRepositories `json:"results"`
	} `json:"data"`
}

type HelmSearchRes struct {
	Code int                             `json:"code"`
	Data []commontypes.ChartListResponse `json:"data"`
}

type HelmReleaseRes struct {
	Code int                          `json:"code"`
	Data []commontypes.ReleaseElement `json:"data"`
}
