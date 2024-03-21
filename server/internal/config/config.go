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
	Consul ConsulConf
}

type ConsulConf struct {
	Address string
}

func NewConfig(configFile, logLevel, consulAddr, listenOn, metricsListenOn *string) Config {
	var c = Config{}
	if *configFile != "" {
		conf.Load(*configFile, &c)
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
	}
	if *logLevel != "" {
		c.Log.Level = *logLevel
	}
	if *consulAddr != "" {
		c.Consul.Address = *consulAddr
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

	if c.Consul.Address == "" {
		logx.Errorf("Consul address must be specified.")
		os.Exit(0)
	}
	return c
}
