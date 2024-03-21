package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf
	Consul           consul.Conf
	Kubeconfig       string
	RegisterEndpoint string
}

func NewConfig(configFile, logLevel, consulHost, consulKey, kubeconfig, listenOn, metricsListenOn, RegisterEndpoint *string) Config {
	var c = Config{}
	if *configFile != "" {
		conf.Load(*configFile, &c)
	} else {
		c.Name = "LizardAgent"
		c.ListenOn = "0.0.0.0:5017"
		c.Timeout = 60000
		c.Log.Encoding = "plain"
		c.Log.Level = "info"
		c.Prometheus.Host = "0.0.0.0"
		c.Prometheus.Port = 15017
		c.Prometheus.Path = "/metrics"
		meta := make(map[string]string)
		c.Consul = consul.Conf{
			Meta: meta,
			TTL:  60,
		}
	}
	if *logLevel != "" {
		c.Log.Level = *logLevel
	}
	if *consulHost != "" {
		c.Consul.Host = *consulHost
	}
	if *consulKey != "" {
		c.Consul.Key = *consulKey
		arr := strings.Split(*consulKey, ".")
		c.Consul.Meta["Protocol"] = "grpc"
		c.Consul.Meta["Service"] = arr[0]
		c.Consul.Meta["Namespace"] = arr[1]
		c.Consul.Meta["Cluster"] = arr[2]
	}
	if *kubeconfig != "" {
		c.Kubeconfig = *kubeconfig
	}
	if *listenOn != "" {
		c.ListenOn = *listenOn
	}
	if *metricsListenOn != "" {
		arr := strings.Split(*metricsListenOn, ":")
		c.Prometheus.Host = arr[0]
		port, _ := strconv.Atoi(arr[1])
		c.Prometheus.Port = port
	}
	if *RegisterEndpoint != "" {
		c.RegisterEndpoint = *RegisterEndpoint
	}
	if c.RegisterEndpoint == "" {
		c.RegisterEndpoint = c.ListenOn
	}
	if c.Consul.Host == "" || c.Consul.Key == "" {
		logx.Errorf("Consul host and key must be specified.")
		os.Exit(0)
	}
	return c
}
