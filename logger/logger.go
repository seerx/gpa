package logger

import (
	"github.com/seerx/logo"
	"github.com/seerx/logo/log"
)

type GpaLogger interface {
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(err error, v ...interface{})
	Errorf(err error, format string, v ...interface{})

	SetLogSQL(log bool)
	IsLogSQL() bool
}

type logger struct {
	logo.Logger
	logSQL bool
}

func (l *logger) SetLogSQL(log bool)                { l.logSQL = log }
func (l *logger) IsLogSQL() bool                    { return l.logSQL }
func (l *logger) Error(err error, v ...interface{}) { l.WithError(err).Error(v...) }
func (l *logger) Errorf(err error, format string, v ...interface{}) {
	l.WithError(err).Errorf(format, v...)
}

var logInstance *logger

func init() {
	logInstance = &logger{
		Logger: log.GetDefaultLogger(),
		logSQL: false,
	}
}

func GetLogger() *logger {
	return logInstance
}
