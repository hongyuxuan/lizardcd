package types

import (
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
