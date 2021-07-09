package logs

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/seerx/gpa/rt/log"
)

type logger struct {
	// showSQL       bool
	showCallStack bool
	fmt           log.Formatter
	logLevel      log.LogLevel
	sqlLog        log.LogPrinter
	// errLog        log.LogPrinter
	// debugLog      log.LogPrinter
	// infoLog       log.LogPrinter
	// warnLog       log.LogPrinter
	loggers   []log.LogPrinter
	mutex     sync.Mutex
	entryPool sync.Pool
}

var (
	defautLogger *logger
	once         sync.Once
)

func GetDefaultLogger() *logger {
	once.Do(func() {
		loggers := []log.LogPrinter{
			log.NewStdLogPrinter(io.Discard),
			log.NewStdLogPrinter(os.Stdout),
			log.NewStdLogPrinter(os.Stdout),
			log.NewStdLogPrinter(os.Stdout),
		}
		defautLogger = &logger{
			fmt:      log.TextFormatter,
			logLevel: log.LevelInfo,
			loggers:  loggers,
			sqlLog:   log.NewStdLogPrinter(io.Discard),
			// debugLog: loggers[0],
			// infoLog:  loggers[1],
			// warnLog:  loggers[2],
			// errLog:   loggers[3],
		}
	})
	return defautLogger
}

func (l *logger) newEntry() *log.Entry {
	entry, ok := l.entryPool.Get().(*log.Entry)
	if ok {
		return entry
	}
	return &log.Entry{
		Logger: l,
	}
}

func (l *logger) releaseEntry(entry *log.Entry) {
	entry.Data = nil
	entry.Err = nil
	l.entryPool.Put(entry)
}

func (l *logger) WithData(data interface{}) *log.Entry {
	e := l.newEntry()
	e.Data = e
	return e
}
func (l *logger) WithError(err error) *log.Entry {
	e := l.newEntry()
	e.Err = err
	return e
}

func (l *logger) log(level log.LogLevel, e *log.Entry) {
	e.Time = time.Now()
	e.Level = level

	msg, err := l.fmt(e)
	if err != nil {
		l.WithError(err).Error("log format error")
		return
	}

	l.loggers[level].Print(msg)
}

func (l *logger) LogSQL(v ...interface{}) {
	e := l.newEntry()
	e.Message = fmt.Sprint(v...)
	e.Time = time.Now()
	e.Level = log.LevelInfo

	msg, err := l.fmt(e)
	if err != nil {
		l.WithError(err).Error("log format error")
		return
	}

	l.sqlLog.Print(msg)
}
func (l *logger) Debug(v ...interface{}) {
	e := l.newEntry()
	e.Message = fmt.Sprint(v...)
	l.log(log.LevelDebug, e)
}
func (l *logger) Debugf(format string, v ...interface{}) {
	e := l.newEntry()
	e.Message = fmt.Sprintf(format, v...)
	l.log(log.LevelDebug, e)
}
func (l *logger) Info(v ...interface{}) {
	e := l.newEntry()
	e.Message = fmt.Sprint(v...)
	l.log(log.LevelInfo, e)
}
func (l *logger) Infof(format string, v ...interface{}) {
	e := l.newEntry()
	e.Message = fmt.Sprintf(format, v...)
	l.log(log.LevelInfo, e)
}
func (l *logger) Warn(v ...interface{}) {
	e := l.newEntry()
	e.Message = fmt.Sprint(v...)
	l.log(log.LevelWarn, e)
}
func (l *logger) Warnf(format string, v ...interface{}) {
	e := l.newEntry()
	e.Message = fmt.Sprintf(format, v...)
	l.log(log.LevelWarn, e)
}
func (l *logger) Error(v ...interface{}) {
	e := l.newEntry()
	e.Message = fmt.Sprint(v...)
	l.log(log.LevelError, e)
}
func (l *logger) Errorf(format string, v ...interface{}) {
	e := l.newEntry()
	e.Message = fmt.Sprintf(format, v...)
	l.log(log.LevelError, e)
}

func (l *logger) ShowSQL(print bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	if print {
		l.sqlLog.SetOutput(os.Stdout)
	} else {
		l.sqlLog.SetOutput(io.Discard)
	}
}
func (l *logger) SetLevel(level log.LogLevel) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.logLevel = level
	for n, logger := range l.loggers {
		if int(level) >= n {
			logger.SetOutput(os.Stdout)
		} else {
			logger.SetOutput(io.Discard)
		}
	}
}
func (l *logger) GetLevel() log.LogLevel         { return l.logLevel }
func (l *logger) ShowCallStacks(show bool)       { l.showCallStack = show }
func (l *logger) SetFormatter(fmt log.Formatter) { l.fmt = fmt }
