package utils

import (
	"context"
	"strings"
)

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

func GetPayload(ctx context.Context) (username, role, tenant string, namespaces []string) {
	payload := ctx.Value("payloads").(map[string]interface{})
	username = payload["username"].(string)
	role = payload["role"].(string)
	tenant = payload["tenant"].(string)
	namespaceStr := payload["namespace"].(string)
	namespaces = strings.Split(namespaceStr, ",")
	return
}
