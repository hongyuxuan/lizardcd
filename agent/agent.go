package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/internal/config"
	"github.com/hongyuxuan/lizardcd/agent/internal/server"
	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"github.com/zeromicro/zero-contrib/zrpc/registry/nacos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	configFile       = kingpin.Flag("config", "config file").Short('f').Default("").String()
	logLevel         = kingpin.Flag("log.level", "Log level.").Default("").String()
	consulHost       = kingpin.Flag("consul-host", "Consul hosts.").Default("").String()
	etcdHost         = kingpin.Flag("etcd-host", "Etcd hosts.").Default("").String()
	nacosHost        = kingpin.Flag("nacos-host", "Nacos hosts.").Default("").String()
	nacosNamespaceId = kingpin.Flag("nacos-namespace-id", "Nacos namespaceId.").Default("").String()
	nacosUsername    = kingpin.Flag("nacos-username", "Nacos username.").Default("").String()
	nacosPassword    = kingpin.Flag("nacos-password", "Nacos password.").Default("").String()
	nacosGroup       = kingpin.Flag("nacos-group", "Nacos group.").Default("").String()
	serviceKey       = kingpin.Flag("service-key", "Service key for registry. Format must be: lizardcd-agent.<namespace>.<cluster>").Default("").String()
	servicePrefix    = kingpin.Flag("service-prefix", "Prefix of service key for registry. Can be empty").Default("").String()
	kubeconfig       = kingpin.Flag("kubeconfig", "Kubeconfig file, must be specified when agent is out-of-k8s deployed").Default("").String()
	listenOn         = kingpin.Flag("grpc-addr", "Grpc listen address.").Default("").String()
	metricsListenOn  = kingpin.Flag("metrics-addr", "Prometheus metrics listen address.").Default("").String()

	/* print app version */
	AppVersion = "unknown"
	GoVersion  = "unknown"
	BuildTime  = "unknown"
	OsArch     = "unknown"
	Author     = "unknown"
)

func main() {
	// Parse flags
	kingpin.Version(printVersion())
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logx.MustSetup(logx.LogConf{Encoding: "plain"})

	c := config.NewConfig(
		configFile,
		logLevel,
		consulHost,
		etcdHost,
		nacosHost,
		nacosNamespaceId,
		nacosUsername,
		nacosPassword,
		nacosGroup,
		serviceKey,
		servicePrefix,
		kubeconfig,
		listenOn,
		metricsListenOn)

	logx.DisableStat()
	logx.MustSetup(c.Log)

	ctx := svc.NewServiceContext(c)
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		agent.RegisterLizardAgentServer(grpcServer, server.NewLizardAgentServer(ctx))
		logx.Infof("Lizardcd-agent: %s register to etcd success", c.Etcd.Key)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// register service to consul
	if c.Consul.Host != "" {
		go registerConsul(c)
	}

	// register service to nacos
	if c.Nacos.Host != "" {
		go registerNacos(c)
	}

	logx.Infof("Starting rpc server at %s...", c.ListenOn)
	s.Start()
}

func printVersion() string {
	return fmt.Sprintf("App version: %s\nGo version:  %s\nBuild Time:  %s\nOS/Arch:     %s\nAuthor:      %s\n", AppVersion, GoVersion, BuildTime, OsArch, Author)
}

func registerConsul(c config.Config) {
	var podIp string
	if podIp = os.Getenv("POD_IP"); podIp == "" {
		podIp = c.ListenOn
	} else {
		port := strings.Split(c.ListenOn, ":")[1]
		podIp += ":" + port
	}
	var err error
	for {
		if err = consul.RegisterService(podIp, c.Consul); err != nil {
			logx.Error(err)
			time.Sleep(time.Duration(5) * time.Second)
		} else {
			logx.Infof("Lizardcd-agent: %s register to consul success", c.Consul.Key)
			break
		}
	}
}

func registerNacos(c config.Config) {
	var podIp string
	if podIp = os.Getenv("POD_IP"); podIp == "" {
		podIp = c.ListenOn
	} else {
		port := strings.Split(c.ListenOn, ":")[1]
		podIp += ":" + port
	}
	fields := strings.Split(c.Nacos.Host, ":")
	host := fields[0]
	port, _ := strconv.ParseUint(fields[1], 10, 64)
	var err error
	for {
		opts := nacos.NewNacosConfig(c.Nacos.Key, podIp, []constant.ServerConfig{
			*constant.NewServerConfig(host, port),
		}, &constant.ClientConfig{
			NamespaceId: c.Nacos.NamespaceId,
			LogLevel:    c.Log.Level,
			LogDir:      "./logs",
			Username:    c.Nacos.Username,
			Password:    c.Nacos.Password,
		}, nacos.WithGroup(c.Nacos.Group), nacos.WithMetadata(c.Nacos.Meta))
		if err = nacos.RegisterService(opts); err != nil {
			logx.Errorf("Failed to connect to nacos: %v", err)
			time.Sleep(time.Duration(5) * time.Second)
		} else {
			logx.Infof("Lizardcd-agent: %s register to nacos success", c.Nacos.Key)
			break
		}
	}
}
