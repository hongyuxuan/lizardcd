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
	Consul                 consul.Conf `json:",optional"`
	Nacos                  NacosConf   `json:",optional"`
	Kubeconfig             string      `json:",optional"`
	ServicePrefix          string      `json:",optional"`
	KubernetesSecretPrefix string      `json:",optional"`
}

type NacosConf struct {
	Host        string
	Key         string
	NamespaceId string
	Group       string
	Username    string
	Password    string
	Meta        map[string]string
}

func NewConfig(configFile, logLevel, consulHost, etcdHost, nacosHost, nacosNamespaceId, nacosUsername, nacosPassword, nacosGroup,
	serviceKey, servicePrefix, kubeconfig, listenOn, metricsListenOn *string) Config {
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
		c.Etcd = discov.EtcdConf{}
		c.Consul = consul.Conf{
			Meta: meta,
			TTL:  60,
		}
		c.Nacos = NacosConf{
			Meta: meta,
		}
	}
	if *logLevel != "" {
		c.Log.Level = *logLevel
	}
	if *servicePrefix != "" {
		c.ServicePrefix = *servicePrefix
	}
	if *consulHost != "" {
		c.Consul.Host = *consulHost
		if *serviceKey != "" {
			meta, err := utils.GetServiceMata(*servicePrefix, *serviceKey)
			if err != nil {
				logx.Error(err)
				os.Exit(0)
			}
			c.Consul.Key = *serviceKey
			c.Consul.Meta = meta
		}
	}
	if *etcdHost != "" {
		c.Etcd.Hosts = strings.Split(*etcdHost, ",")
		if *serviceKey != "" {
			c.Etcd.Key = *serviceKey
		}
	}
	if *nacosHost != "" {
		c.Nacos.Host = *nacosHost
		if *serviceKey != "" {
			c.Nacos.Key = *serviceKey
		}
	}
	if *nacosNamespaceId != "" {
		c.Nacos.NamespaceId = *nacosNamespaceId
	}
	if *nacosUsername != "" {
		c.Nacos.Username = *nacosUsername
	}
	if *nacosPassword != "" {
		c.Nacos.Password = *nacosPassword
	}
	if *nacosGroup != "" {
		c.Nacos.Group = *nacosGroup
	}
	c.Consul.Key = c.ServicePrefix + c.Consul.Key
	c.Etcd.Key = c.ServicePrefix + c.Etcd.Key
	c.Nacos.Key = c.ServicePrefix + c.Nacos.Key
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
	if c.KubernetesSecretPrefix == "" {
		c.KubernetesSecretPrefix = "default-token" // default token prefix
	}
	if len(c.Etcd.Hosts) == 0 && c.Consul.Host == "" && c.Nacos.Host == "" {
		logx.Errorf("Either etcd host, consul host or nacos host must be specified.")
		os.Exit(0)
	}
	logx.Infof("Using config: %+v", c)
	return c
}
