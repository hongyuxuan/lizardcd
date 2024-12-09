package utils

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

type MyFormatter struct{}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	timestamp := entry.Time.Format("2006/01/02 15:04:05.000")
	logcontent := fmt.Sprintf("%s [%s] %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	b.WriteString(logcontent)
	return b.Bytes(), nil
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
