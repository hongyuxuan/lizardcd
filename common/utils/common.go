package utils

import "strings"

func GetLizardAgentKey(key []byte) string {
	arr := strings.Split(string(key), "/")
	uid := arr[len(arr)-1]
	return strings.TrimSuffix(string(key), "/"+uid)
}

func GetServiceMata(key string) map[string]string {
	arr := strings.Split(key, ".")
	return map[string]string{
		"Protocol":  "grpc",
		"Service":   arr[0],
		"Namespace": arr[1],
		"Cluster":   arr[2],
	}
}
