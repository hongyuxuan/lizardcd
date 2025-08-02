package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var tracer = otel.Tracer("gorm/tidb")

type Tidb struct {
	*gorm.DB
}

func NewTidb(host, username, password, database, level string, port int) *gorm.DB {
	InitLogger(level)
	loglevel := logger.Silent
	if level == "debug" {
		loglevel = logger.Info
	}
	tidbLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second * 1,
			LogLevel:      loglevel,
		},
	)
	dburl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local", username, password, host, port, database)
	tidb, err := gorm.Open(mysql.Open(dburl), &gorm.Config{
		Logger: tidbLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		Log.Fatalf("Failed to connect tidb: %v", err)
	}
	Log.Infof("Init tidb connection %s success", dburl)

	// callback
	tidb.Callback().Create().Before("gorm:before_create").Register("callback_before", tracingBefore)
	tidb.Callback().Query().Before("gorm:before_query").Register("callback_before", tracingBefore)
	tidb.Callback().Update().Before("gorm:before_update").Register("callback_before", tracingBefore)
	tidb.Callback().Delete().Before("gorm:before_delete").Register("callback_before", tracingBefore)

	tidb.Callback().Create().After("gorm:after_create").Register("callback_after", tracingAfter)
	tidb.Callback().Query().After("gorm:after_query").Register("callback_after", tracingAfter)
	tidb.Callback().Update().After("gorm:after_update").Register("callback_after", tracingAfter)
	tidb.Callback().Delete().After("gorm:after_delete").Register("callback_after", tracingAfter)
	return tidb
}
