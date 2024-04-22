package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils"
)

var tracer = otel.Tracer("gorm/sqlite")

type SQLite struct {
	*gorm.DB
}

func NewSQLite(dbfile, level string) *gorm.DB {
	InitLogger(level)
	loglevel := logger.Silent
	if level == "debug" {
		loglevel = logger.Info
	}
	sqliteLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second * 1,
			LogLevel:      loglevel,
		},
	)
	sqlite, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{
		Logger: sqliteLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		logx.Errorf("failed to connect sqlite: %v", err)
		os.Exit(0)
	}
	logx.Infof("open sqlite file %s success", dbfile)

	// callback
	sqlite.Callback().Create().Before("gorm:before_create").Register("callback_before", tracingBefore)
	sqlite.Callback().Query().Before("gorm:before_query").Register("callback_before", tracingBefore)
	sqlite.Callback().Update().Before("gorm:before_update").Register("callback_before", tracingBefore)
	sqlite.Callback().Delete().Before("gorm:before_delete").Register("callback_before", tracingBefore)

	sqlite.Callback().Create().After("gorm:after_create").Register("callback_after", tracingAfter)
	sqlite.Callback().Query().After("gorm:after_query").Register("callback_after", tracingAfter)
	sqlite.Callback().Update().After("gorm:after_update").Register("callback_after", tracingAfter)
	sqlite.Callback().Delete().After("gorm:after_delete").Register("callback_after", tracingAfter)
	return sqlite
}

func SetTx(tx *gorm.DB, count *int64, req *commontypes.GetDataReq) {
	if req.Search != "" {
		searchStmt := strings.Split(req.Search, "==")
		tx.Where(fmt.Sprintf("%s LIKE ?", searchStmt[0]), "%"+searchStmt[1]+"%")
	}
	if req.Filter != "" {
		for _, filter := range strings.Split(req.Filter, ",") {
			filterStmt := strings.Split(filter, "==")
			tx.Where(fmt.Sprintf("%s = ?", filterStmt[0]), filterStmt[1])
		}
	}
	if req.Range != "" { // &range=init_at==2024-02-06 16:58:30,2024-02-06 17:58:30
		rangeS := strings.Split(req.Range, "==")
		rangeKey := rangeS[0]
		rangeR := strings.Split(rangeS[1], ",")
		tx.Where(fmt.Sprintf("%s BETWEEN ? AND ?", rangeKey), rangeR[0], rangeR[1])
	}
	tx.Count(count)
	tx = tx.Limit(req.Size).Offset((req.Page - 1) * req.Size)
	if req.Sort != "" {
		tx = tx.Order(req.Sort)
	}
}

func tracingBefore(db *gorm.DB) {
	db.InstanceSet("startAt", time.Now())
}

func tracingAfter(db *gorm.DB) {
	stmt := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	ctx := db.Statement.Context
	spanName, ok := ctx.Value("SpanName").(string)
	if !ok {
		spanName = "TiDB"
	}
	_, span := tracer.Start(ctx, spanName)
	defer span.End()
	if startAt, ok := db.InstanceGet("startAt"); ok {
		span.SetAttributes(
			attribute.Int64("database.cost", time.Since(startAt.(time.Time)).Milliseconds()),
		)
	}
	// var res []map[string]interface{}
	// db.Scan(&res)
	// resJson, _ := json.MarshalIndent(res, "", "  ")
	span.SetAttributes(
		attribute.String("database.statement", stmt),
		attribute.Int64("database.rows_affected", db.RowsAffected),
		// attribute.String("database.result", string(resJson)),
		attribute.String("database.stack", utils.FileWithLineNum()),
	)
	if db.Error != nil {
		span.RecordError(db.Error)
	}
}
