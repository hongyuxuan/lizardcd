package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	Consul        ConsulConf `json:",optional"`
	Nacos         NacosConf  `json:",optional"`
	Etcd          EtcdConf   `json:",optional"`
	ServicePrefix string     `json:",optional"`
	Sqlite        string
}

type NacosConf struct {
	Address     string
	NamespaceId string
	Group       string
	Username    string
	Password    string
}

type EtcdConf struct {
	Address string
}

type ConsulConf struct {
	Address string
}

func NewConfig(configFile, logLevel, consulAddr, etcdAddr, nacosAddr, nacosNamespaceId, nacosUsername, nacosPassword, nacosGroup, servicePrefix, listenOn, metricsListenOn, dbfile, accessSecret *string, accessExpire *int64) Config {
	var c = Config{}
	if *configFile != "" {
		conf.MustLoad(*configFile, &c)
	} else {
		c.Name = "LizardServer"
		c.Host = "0.0.0.0"
		c.Port = 5117
		c.Timeout = 60000
		c.Log.Encoding = "plain"
		c.Log.Level = "info"
		c.Prometheus.Host = "0.0.0.0"
		c.Prometheus.Port = 15117
		c.Prometheus.Path = "/metrics"
		c.Consul = ConsulConf{}
		c.Etcd = EtcdConf{}
		c.Nacos = NacosConf{}
		c.Sqlite = "./lizardcd.db"
		c.Auth.AccessSecret = "wLnOk8keh/WO5u7lX8H1dB1/mcuHvnI/jfWCMXMPg9o="
		c.Auth.AccessExpire = 86400
	}
	if *logLevel != "" {
		c.Log.Level = *logLevel
	}
	if *etcdAddr != "" {
		c.Etcd.Address = *etcdAddr
	}
	if *consulAddr != "" {
		c.Consul.Address = *consulAddr
	}
	if *nacosAddr != "" {
		c.Nacos.Address = *nacosAddr
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
	if *listenOn != "" {
		arr := strings.Split(*listenOn, ":")
		c.Host = arr[0]
		port, _ := strconv.Atoi(arr[1])
		c.Port = port
	}
	if *metricsListenOn != "" {
		arr := strings.Split(*metricsListenOn, ":")
		c.Prometheus.Host = arr[0]
		port, _ := strconv.Atoi(arr[1])
		c.Prometheus.Port = port
	}
	if *servicePrefix != "" {
		c.ServicePrefix = *servicePrefix
	}
	if *dbfile != "" {
		c.Sqlite = *dbfile
	}
	if *accessSecret != "" {
		c.Auth.AccessSecret = *accessSecret
	}
	if *accessExpire != 0 {
		c.Auth.AccessExpire = *accessExpire
	}

	if c.Etcd.Address == "" && c.Consul.Address == "" && c.Nacos.Address == "" {
		logx.Errorf("Either etcd, consul or nacos address must be specified.")
		os.Exit(0)
	}
	return c
}
