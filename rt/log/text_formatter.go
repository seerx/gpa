package log

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
)

// TextFormatter 自定义日志格式化输出
type TextFormatter struct {
	timeFormat string
}

const (
	formatWithCaller = "[%s] %s %s\n%s\n%s:%d"
	format           = "[%s] %s %s"
)

// Format 格式化日志信息
func (t *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var body string
	if entry.HasCaller() {
		body = fmt.Sprintf(formatWithCaller,
			entry.Level.String(),
			entry.Time.Format(t.timeFormat),
			entry.Message,
			entry.Caller.Func.Name(),
			entry.Caller.File,
			entry.Caller.Line)
	} else {
		body = fmt.Sprintf(format,
			entry.Level.String(),
			entry.Time.Format(t.timeFormat),
			entry.Message)
	}
	b := &bytes.Buffer{}
	b.WriteString(body)

	for k, v := range entry.Data {
		b.WriteString(fmt.Sprintf("\n%s:%v", k, v))
	}

	b.WriteString("\n")
	return b.Bytes(), nil
}
