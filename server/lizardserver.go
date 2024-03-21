package main

import (
	"fmt"
	"net/http"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/hongyuxuan/lizardcd/server/internal/config"
	"github.com/hongyuxuan/lizardcd/server/internal/handler"
	"github.com/hongyuxuan/lizardcd/server/internal/svc"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"

	_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

var (
	configFile      = kingpin.Flag("config", "config file").Short('f').Default("").String()
	consulAddr      = kingpin.Flag("consul-addr", "Consul address.").Default("").String()
	logLevel        = kingpin.Flag("log.level", "Log level.").Default("").String()
	listenOn        = kingpin.Flag("http-addr", "HTTP listen address.").Default("").String()
	metricsListenOn = kingpin.Flag("metrics-addr", "Prometheus metrics listen address.").Default("").String()

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

	c := config.NewConfig(configFile, logLevel, consulAddr, listenOn, metricsListenOn)
	logx.DisableStat()
	logx.MustSetup(c.Log)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case *errorx.LizardcdError:
			return e.Code, e.GetData()
		default:
			return http.StatusInternalServerError, errorx.HttpErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
	})

	go handler.StartConsulWatch(ctx)

	logx.Infof("Starting server at %s:%d...", c.Host, c.Port)
	server.Start()
}

func printVersion() string {
	return fmt.Sprintf("App version: %s\nGo version:  %s\nBuild Time:  %s\nOS/Arch:     %s\nAuthor:      %s\n", AppVersion, GoVersion, BuildTime, OsArch, Author)
}
