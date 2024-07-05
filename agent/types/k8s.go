package types

type ScaleReq struct {
	Workloads []ScaleWorkload `json:"workloads"`
}

type ScaleWorkload struct {
	Name     string `json:"name"`
	Replicas int    `json:"replicas"`
	Disabled bool   `json:"disabled,optional"`
}
