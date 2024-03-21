package main

import (
	"fmt"
	"time"

	"github.com/hongyuxuan/lizardcd/agent/internal/config"
	"github.com/hongyuxuan/lizardcd/agent/internal/server"
	"github.com/hongyuxuan/lizardcd/agent/internal/svc"
	"github.com/hongyuxuan/lizardcd/agent/types/agent"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	configFile       = kingpin.Flag("config", "config file").Short('f').Default("").String()
	logLevel         = kingpin.Flag("log.level", "Log level.").Default("").String()
	consulHost       = kingpin.Flag("consul.host", "Consul host.").Default("").String()
	consulKey        = kingpin.Flag("consul.key", "Consul service key. Format must be: lizardcd-agent.<namespace>.<cluster>").Default("").String()
	kubeconfig       = kingpin.Flag("kubeconfig", "Kubeconfig file, must be specified when agent is out-of-k8s deployed").Default("").String()
	listenOn         = kingpin.Flag("grpc-addr", "Grpc listen address.").Default("").String()
	metricsListenOn  = kingpin.Flag("metrics-addr", "Prometheus metrics listen address.").Default("").String()
	RegisterEndpoint = kingpin.Flag("register.endpoint", "Service endpoint to register to consul.").Default("").String()

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

	c := config.NewConfig(configFile, logLevel, consulHost, consulKey, kubeconfig, listenOn, metricsListenOn, RegisterEndpoint)

	logx.DisableStat()
	logx.MustSetup(c.Log)

	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		agent.RegisterLizardAgentServer(grpcServer, server.NewLizardAgentServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// register service to consul
	go registerConsul(c)

	logx.Infof("Starting rpc server at %s...", c.ListenOn)
	s.Start()
}

func printVersion() string {
	return fmt.Sprintf("App version: %s\nGo version:  %s\nBuild Time:  %s\nOS/Arch:     %s\nAuthor:      %s\n", AppVersion, GoVersion, BuildTime, OsArch, Author)
}

func registerConsul(c config.Config) {
	var err error
	for {
		if err = consul.RegisterService(c.RegisterEndpoint, c.Consul); err != nil {
			logx.Error(err)
			time.Sleep(time.Duration(5) * time.Second)
		} else {
			logx.Infof("Agent register to consul success")
			break
		}
	}
}
