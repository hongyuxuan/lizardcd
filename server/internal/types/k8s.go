package types

type PatchYamlReq struct {
	Cluster   string `path:"cluster"`
	Namespace string `path:"namespace"`
}
