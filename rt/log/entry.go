package log

import (
	"encoding/json"
	"fmt"
	"time"
)

type Entry struct {
	Logger  `json:"-"`
	Time    time.Time   `json:"time"`  // 转换
	Level   LogLevel    `json:"level"` // 转换
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Err     error       `json:"error"`
}

func (e *Entry) WithData(data interface{}) *Entry {
	e.Data = data
	return e
}
func (e *Entry) WithError(err error) *Entry {
	e.Err = err
	return e
}

type Formatter func(*Entry) ([]byte, error)

func TextFormatter(e *Entry) ([]byte, error) {
	data := ""
	if e.Data != nil {
		jsonData, err := json.MarshalIndent(e.Data, "", "  ")
		if err != nil {
			return nil, err
		}
		data = string(jsonData)
	}
	return []byte(fmt.Sprintf("[%5s] %s %s\n%s", GetLevelString(e.Level), e.Time.Format(time.RFC3339), e.Message, data)), nil
}

func JSONFormatter(e *Entry) ([]byte, error) {
	// format = "[%s] %s %s"
	// data := ""
	// if e.Data != nil {
	return json.MarshalIndent(e, "", "  ")
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	// data = string(jsonData)
	// }
	// return []byte(fmt.Sprintf("[%5s] %s %s\n%s", e.Level, e.Time, e.Message, data)), nil
}
