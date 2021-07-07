package dbutil

import (
	"encoding/json"
	"errors"
)

func Struct2String(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ParseStruct(src interface{}, bean interface{}) error {
	switch val := src.(type) {
	case string:
		return json.Unmarshal([]byte(val), bean)
	case []byte:
		return json.Unmarshal(val, bean)
	}
	return errors.New("not support data type to convert to struct")
}
