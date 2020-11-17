package logger

import (
	"fmt"

	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logPath string
)

// SetLogPath prepares the path of the logfile
func SetLogPath(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	logPath = path
}

// GetZapLogger get zap logger
func GetZapLogger(logfile string) *zap.Logger {

	logfile = strings.ReplaceAll(logfile, "/", "_")
	logger, exists := zapLoggersMap[logfile]

	if !exists {

		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%s.log", logPath, logfile),
			MaxSize:    5, // megabytes
			MaxBackups: 10,
			MaxAge:     30, // days
		})
		coreConfig := zap.NewProductionEncoderConfig()
		coreConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(coreConfig),
			w,
			zap.InfoLevel,
		)

		zapLoggersMap[logfile] = zap.New(core)

		return zapLoggersMap[logfile]
	}
	return logger
}

// ##################################################################
// ##################################################################
// The following funcs are added to meet logging requirements of the
// billing and other new services

// ApplicationLog application log
type ApplicationLog struct {
	zapLogger *zap.Logger
	fields    []zapcore.Field
	logLevel  zapcore.Level
}

// AddOne adds one key value
func (lg *ApplicationLog) AddOne(key string, value interface{}) *ApplicationLog {
	switch value.(type) {
	case int:
		lg.fields = append(lg.fields, zap.Int(key, value.(int)))
	case float64:
		lg.fields = append(lg.fields, zap.Float64(key, value.(float64)))
	case string:
		lg.fields = append(lg.fields, zap.String(key, value.(string)))
	case time.Duration:
		lg.fields = append(lg.fields, zap.Duration(key, value.(time.Duration)))
	default:
		lg.fields = append(lg.fields, zap.Any(key, value))
	}

	return lg
}

// Append appends prepared log items to log slice
func (lg *ApplicationLog) Append(fields ...zapcore.Field) *ApplicationLog {

	lg.fields = append(lg.fields, fields...)
	return lg
}

// Commit appends to log
func (lg *ApplicationLog) Commit(message string) {

	defer func() {
		lg.zapLogger.Sync()
		lg.fields = nil
	}()

	switch lg.logLevel {

	case zapcore.InfoLevel:
		lg.zapLogger.Info(message, lg.fields...)

	case zap.PanicLevel:
		lg.zapLogger.Panic(message, lg.fields...)

	case zap.FatalLevel:
		lg.zapLogger.Fatal(message, lg.fields...)

	case zap.WarnLevel:
		lg.zapLogger.Warn(message, lg.fields...)

	default:
		lg.zapLogger.Warn(message, lg.fields...)
	}

}

// Level sets log level
func (lg *ApplicationLog) Level(level zapcore.Level) {
	lg.logLevel = level
}

// AddDevelopementInfo adds debug info
func (lg *ApplicationLog) AddDevelopementInfo() {

	pc, _, line, ok := runtime.Caller(1) // pc, file, line, ok
	details := runtime.FuncForPC(pc)     // all we want to learn about the process
	var callerFunction string
	if ok && details != nil {
		callerFunction = details.Name()
	}
	runtimeInfo := make(map[string]interface{})

	runtimeInfo["line"] = line
	runtimeInfo["caller"] = callerFunction
	lg.AddOne("developer_info", runtimeInfo)
}

// AddDatabaseConnectionPoolInfo reflects info about databaase connection
// func (lg *ApplicationLog) AddDatabaseConnectionPoolInfo(db *model.DBFunc) *ApplicationLog {
// 	lg.AddOne("database", db.GetSqlConnection().Stats())
// 	return lg
// }

// AddHTTPRequestInfo adds details about http
func (lg *ApplicationLog) AddHTTPRequestInfo(c *gin.Context) *ApplicationLog {

	lg.Append(zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("ip", c.ClientIP()),
		zap.String("user-agent", c.Request.UserAgent()))

	return lg
}

// Prepare preps log
func Prepare(zapLogger *zap.Logger) *ApplicationLog {
	return &ApplicationLog{zapLogger: zapLogger}
}
