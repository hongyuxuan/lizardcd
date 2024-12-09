package utils

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardcd/common/constant"
	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/samber/lo"
)

func GetLizardAgentKey(key []byte) string {
	arr := strings.Split(string(key), "/")
	uid := arr[len(arr)-1]
	return strings.TrimSuffix(string(key), "/"+uid)
}

func GetServiceMata(prefix, key string) (map[string]string, error) {
	re, _ := regexp.Compile(prefix + "lizardcd-agent(_vm|)\\.(.+?)\\.(.+)")
	res := re.FindStringSubmatch(key)
	if res == nil {
		return nil, errorx.NewDefaultError("No match of \"%s\" to <ServicePrefix>.lizardcd-agent(_vm).<namespace>.<cluster>", key)
	}
	return map[string]string{
		"Protocol":  "grpc",
		"Service":   prefix + "lizardcd-agent" + res[1],
		"Namespace": res[2],
		"Cluster":   res[3],
	}, nil
}

func GetTarget(prefix, key string, namespaces []string, role string) (target string, err error) {
	re, _ := regexp.Compile(prefix + "lizardcd-agent_vm\\.(.+?)\\.(.+)")
	res := re.FindStringSubmatch(key)
	if res == nil {
		return "", errorx.NewDefaultError("No match of \"%s\" to <ServicePrefix>.lizardcd-agent_vm.<ipaddress>", key)
	}
	if _, ok := lo.Find(namespaces, func(n string) bool {
		return n == res[1]
	}); ok || role == constant.ROLE_ADMIN {
		return res[2], nil
	}
	return "", errorx.NewDefaultError("Permisson denied of \"%s\" for current tenant", key)
}

func GetPayload(ctx context.Context) (username, role string, tenant, namespaces []string) {
	payload := ctx.Value("payloads").(map[string]interface{})
	username = payload["username"].(string)
	role = payload["role"].(string)
	tenantStr := payload["tenant"].(string)
	tenant = strings.Split(tenantStr, ",")
	namespaceStr := payload["namespace"].(string)
	namespaces = strings.Split(namespaceStr, ",")
	return
}

func AnyToString(v interface{}) string {
	switch v.(type) {
	case string:
		return v.(string)
	case bool:
		return strconv.FormatBool(v.(bool))
	case float64:
		return strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	default:
		return fmt.Sprintf("%+v", v)
	}
}

func GetLocalListener(port int) string {
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					return fmt.Sprintf("%s:%d", ipnet.IP.String(), port)
				}
			}
		}
	}
	return ""
}

func GenerateRandomInt(length int) string {
	var charset = "123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
