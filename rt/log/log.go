package log

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

// Log 默认日志
var logger *logrus.Logger

// InitLog 初始化默认日志
func InitLog(b *Builder) {
	if b == nil {
		logger = NewBuilder().
			Build()
	}
	logger = b.Build()
	// runtime.Compiler
	logger.WithField("version", runtime.Version()).Info("Golang " + runtime.Compiler)
}

// Get 获取 logger
func Get() *logrus.Logger {
	return logger
}

// WithFieldExt logger.WithField
func WithFieldExt(key string, value interface{}, request *http.Request) *logrus.Entry {
	return logger.WithFields(logrus.Fields{key: value, "ip": request.RemoteAddr})
}

// WithFieldsExt logger.WithFields
func WithFieldsExt(fields logrus.Fields, request *http.Request) *logrus.Entry {
	fields["ip"] = request.RemoteAddr
	return logger.WithFields(fields)
}

// WithField logger.WithField
func WithField(key string, value interface{}) *logrus.Entry {
	return logger.WithField(key, value)
}

// WithFields logger.WithFields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

// WithError logger.WithError
func WithError(err error) *logrus.Entry {
	return logger.WithError(err)
}

// WithContext Add a context to the log entry.
func WithContext(ctx context.Context) *logrus.Entry {
	return logger.WithContext(ctx)
}

// WithTime Overrides the time of the log entry.
func WithTime(t time.Time) *logrus.Entry {
	return logger.WithTime(t)
}

// WithJSON 打印 JOSN
func WithJSON(obj interface{}) *logrus.Entry {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return logger.WithFields(logrus.Fields{
			"WithJSON-Error": err.Error(),
			"WithJSON-Value": fmt.Sprintf("%v", obj),
		})
	}
	return logger.WithField("JSON", string(data))
}

// Logf logger.Logf
func Logf(level logrus.Level, format string, args ...interface{}) {
	logger.Logf(level, format, args...)
}

// Tracef logger.Tracef
func Tracef(format string, args ...interface{}) {
	logger.Tracef(format, args...)
}

// Debugf logger.Debugf
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Infof logger.Infof
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Printf logger.Printf
func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

// Warnf logger.Warnf
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Warningf logger.Warningf
func Warningf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}

// Errorf logger.Errorf
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Fatalf logger.Fatalf
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

// Panicf logger.Panicf
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

// Log logger.Log
func Log(level logrus.Level, args ...interface{}) {
	logger.Log(level, args...)
}

// Trace logger.Trace
func Trace(args ...interface{}) {
	logger.Trace(args...)
}

// Debug logger.Debug
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Info logger.Info
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Print logger.Print
func Print(args ...interface{}) {
	logger.Print(args...)
}

// Warn logger.Warn
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warning logger.Warning
func Warning(args ...interface{}) {
	logger.Warning(args...)
}

// Error logger.Error
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Fatal logger.Fatal
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Panic logger.Panic
func Panic(args ...interface{}) {
	logger.Panic(args...)
}

// Logln logger.Logln
func Logln(level logrus.Level, args ...interface{}) {
	logger.Logln(level, args...)
}

// Traceln logger.Traceln
func Traceln(args ...interface{}) {
	logger.Traceln(args...)
}

// Debugln logger.Debugln
func Debugln(args ...interface{}) {
	logger.Debugln(args...)
}

// Infoln logger.Infoln
func Infoln(args ...interface{}) {
	logger.Infoln(args...)
}

// Println logger.Println
func Println(args ...interface{}) {
	logger.Println(args...)
}

// Warnln logger.Warnln
func Warnln(args ...interface{}) {
	logger.Warnln(args...)
}

// Warningln logger.Warningln
func Warningln(args ...interface{}) {
	logger.Warningln(args...)
}

// Errorln logger.Errorln
func Errorln(args ...interface{}) {
	logger.Errorln(args...)
}

// Fatalln logger.Fatalln
func Fatalln(args ...interface{}) {
	logger.Fatalln(args...)
}

// Panicln logger.Panicln
func Panicln(args ...interface{}) {
	logger.Panicln(args...)
}

// Exit logger.Exit
// func Exit(code int) {
// 	logger.Exit(code)
// }

//SetNoLock When file is opened with appending mode, it's safe to
//write concurrently to a file (within 4k message on Linux).
//In these cases user can choose to disable the lock.
func SetNoLock() {
	logger.SetNoLock()
}

// SetLevel sets the logger level.
// func SetLevel(level logrus.Level) {
// 	logger.SetLevel(level)
// }

// GetLevel returns the logger level.
// func GetLevel() logrus.Level {
// 	return logger.GetLevel()
// }

// AddHook adds a hook to the logger hooks.
// func AddHook(hook logrus.Hook) {
// 	logger.AddHook(hook)
// }

// IsLevelEnabled checks if the log level of the logger is greater than the level param
// func IsLevelEnabled(level logrus.Level) bool {
// 	return logger.IsLevelEnabled(level)
// }

// SetFormatter sets the logger formatter.
// func SetFormatter(formatter logrus.Formatter) {
// 	logger.SetFormatter(formatter)
// }

// SetOutput sets the logger output.
// func SetOutput(output io.Writer) {
// 	logger.SetOutput(output)
// }

// SetReportCaller logger.SetReportCaller
// func SetReportCaller(reportCaller bool) {
// 	logger.SetReportCaller(reportCaller)
// }

// ReplaceHooks replaces the logger hooks and returns the old ones
// func ReplaceHooks(hooks logrus.LevelHooks) logrus.LevelHooks {
// 	return logger.ReplaceHooks(hooks)
// }
