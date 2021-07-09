package log

import (
	"fmt"
	"io"
)

type LogPrinter interface {
	SetOutput(writer io.Writer)
	Print([]byte)
	// Printf(fmt string, args ...interface{})
}

type StdLogPrinter struct {
	writer io.Writer
}

func NewStdLogPrinter(writer io.Writer) *StdLogPrinter {
	return &StdLogPrinter{
		writer: writer,
	}
}

func (s *StdLogPrinter) SetOutput(writer io.Writer) {
	s.writer = writer
}

// func (s *StdLogPrinter) Printf(fmter string, args ...interface{}) {
// 	fmt.Fprintf(s.writer, fmter, args...)
// }

func (s *StdLogPrinter) Print(message []byte) {
	fmt.Fprint(s.writer, string(message))
}

type Logger interface {
	WithData(data interface{}) *Entry
	WithError(err error) *Entry
	LogSQL(v ...interface{})
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})

	ShowSQL(print bool)
	SetLevel(level LogLevel)
	GetLevel() LogLevel
	ShowCallStacks(show bool)
	SetFormatter(fmt Formatter)
}

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

var levels = []string{
	"debug", "info", "warn", "error",
}

func GetLevelString(level LogLevel) string {
	return levels[level]
}
