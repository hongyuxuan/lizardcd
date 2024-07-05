package types

type WorkloadStatus struct {
	Name string      `json:"name"`
	Pods []PodStatus `json:"pod_status"`
}

type PodStatus struct {
	PodName string `json:"pod_name"`
	Ready   string `json:"ready"`
}
