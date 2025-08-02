package types

type ListResponseData struct {
	Results  interface{} `json:"results"`
	Continue string      `json:"continue,omitempty"`
}

type CommonK8sResource struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	}
}

type WorkloadStatus struct {
	Name     string      `json:"name"`
	Revision string      `json:"revision,omitempty"`
	Pods     []PodStatus `json:"pod_status"`
	Status   string      `json:"status,omitempty"`
	Ready    bool        `json:"ready,omitempty"`
}

type WorkloadImage struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type WorkloadReplica struct {
	Name     string `json:"name"`
	Replicas int32  `json:"replica"`
}

type PodStatus struct {
	PodName string            `json:"pod_name"`
	Ready   string            `json:"ready"`
	Image   map[string]string `json:"image,omitempty"`
}
