package common

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/hongyuxuan/lizardcd/common/utils"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

var LogLevel string
var ConfigFile string
var LizardServer *utils.HttpClient

func InitConfig() {
	utils.InitLogger(LogLevel)

	if strings.Contains(ConfigFile, "~") {
		u, _ := user.Current()
		ConfigFile = strings.ReplaceAll(ConfigFile, "~", u.HomeDir)
	}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(ConfigFile)
	var err error
	var f *os.File
	if _, err = os.Stat(ConfigFile); err != nil && os.IsNotExist(err) {
		if f, err = os.Create(ConfigFile); err != nil {
			utils.Log.Fatalf("failed to create config file \"%s\": %v", ConfigFile, err)
		}
		defer f.Close()
	}
	viper.ReadInConfig()
	serverAddr := viper.GetString("lizardcd.server.url")
	if serverAddr == "" {
		viper.Set("lizardcd.server.url", "http://localhost:5117") // set default lizardcd-server address
		viper.WriteConfig()
		serverAddr = "http://localhost:5117"
	}
	LizardServer = utils.NewHttpClient(otel.Tracer("imroc/req"))
	LizardServer.SetBaseURL(serverAddr)
	if access_token := viper.GetString("lizardcd.auth.access_token"); access_token != "" {
		LizardServer.SetCommonBearerAuthToken(access_token)
	}
	if LogLevel == "debug" {
		LizardServer.EnableDebugLog()
		LizardServer.EnableDumpAll()
	}
	utils.Log.Debugf("init lizardcd-server client %s success", serverAddr)
}

func GetExec() string {
	path, _ := os.Executable()
	_, exec := filepath.Split(path)
	return exec
}

func PrintError(format string, a ...any) {
	format = "\033[0;31;40m" + format + "\033[0m\n"
	fmt.Printf(format, a...)
}

func PrintSuccess(format string, a ...any) {
	format = "\033[0;32;40m" + format + "\033[0m\n"
	fmt.Printf(format, a...)
}
