package types

type ScaleReq struct {
	Workloads []ScaleWorkload `json:"workloads"`
}

type ScaleWorkload struct {
	Name     string `json:"name"`
	Replicas int    `json:"replicas"`
	Disabled bool   `json:"disabled,optional"`
}

type WorkloadStatus struct {
	Name string      `json:"name"`
	Pods []PodStatus `json:"pod_status"`
}

type PodStatus struct {
	PodName string `json:"pod_name"`
	Ready   string `json:"ready"`
}
