package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

type Config struct {
	zrpc.RpcServerConf
	Consul     consul.Conf `json:",optional"`
	Kubeconfig string      `json:",optional"`
}

func NewConfig(configFile, logLevel, consulHost, etcdHost, serviceKey, kubeconfig, listenOn, metricsListenOn *string) Config {
	var c = Config{}
	if *configFile != "" {
		conf.MustLoad(*configFile, &c)
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
		c.Etcd = discov.EtcdConf{}
	}
	if *logLevel != "" {
		c.Log.Level = *logLevel
	}
	if *consulHost != "" {
		c.Consul.Host = *consulHost
		if *serviceKey != "" {
			c.Consul.Key = *serviceKey
			c.Consul.Meta = utils.GetServiceMata(*serviceKey)
		}
	}
	if *etcdHost != "" {
		c.Etcd.Hosts = strings.Split(*etcdHost, ",")
		if *serviceKey != "" {
			c.Etcd.Key = *serviceKey
		}
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
	if c.Consul.Host == "" && len(c.Etcd.Hosts) == 0 {
		logx.Errorf("Either consul host or etcd host must be specified.")
		os.Exit(0)
	}
	return c
}
