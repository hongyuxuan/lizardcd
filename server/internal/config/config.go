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
	Consul ConsulConf `json:",optional"`
	Etcd   EtcdConf   `json:",optional"`
	Sqlite string
}

type ConsulConf struct {
	Address string
}

type EtcdConf struct {
	Address string
}

func NewConfig(configFile, logLevel, consulAddr, etcdAddr, listenOn, metricsListenOn, dbfile, accessSecret *string, accessExpire *int64) Config {
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
		c.Sqlite = "./lizardcd.db"
		c.Auth.AccessSecret = "wLnOk8keh/WO5u7lX8H1dB1/mcuHvnI/jfWCMXMPg9o="
		c.Auth.AccessExpire = 86400
	}
	if *logLevel != "" {
		c.Log.Level = *logLevel
	}
	if *consulAddr != "" {
		c.Consul.Address = *consulAddr
	}
	if *etcdAddr != "" {
		c.Etcd.Address = *etcdAddr
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
	if *dbfile != "" {
		c.Sqlite = *dbfile
	}
	if *accessSecret != "" {
		c.Auth.AccessSecret = *accessSecret
	}
	if *accessExpire != 0 {
		c.Auth.AccessExpire = *accessExpire
	}

	if c.Consul.Address == "" && c.Etcd.Address == "" {
		logx.Errorf("Either consul or etcd address must be specified.")
		os.Exit(0)
	}
	return c
}
