package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

type MyFormatter struct {
	logrus.TextFormatter
}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 获取调用者信息
	var file string
	var line int
	if entry.HasCaller() {
		file = entry.Caller.File // 只取文件名，不包含路径
		line = entry.Caller.Line
	}

	msg := fmt.Sprintf("%s [%s] %s:%d - %s\n",
		entry.Time.Format("2006-01-02 15:04:05.000"),
		strings.ToUpper(entry.Level.String()),
		file,
		line,
		entry.Message)

	return []byte(msg), nil
}

func InitLogger(logLevel string) {
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)
	Log.SetReportCaller(true)
	if logLevel == "debug" {
		Log.SetLevel(logrus.DebugLevel)
	} else if logLevel == "warning" {
		Log.SetLevel(logrus.WarnLevel)
	} else if logLevel == "error" {
		Log.SetLevel(logrus.ErrorLevel)
	} else if logLevel == "fatal" {
		Log.SetLevel(logrus.FatalLevel)
	}
	Log.SetFormatter(&MyFormatter{})
}
