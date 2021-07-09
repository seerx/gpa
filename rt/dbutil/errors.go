package dbutil

import "fmt"

// ErrParamIsEmpty SQL 的 in 参数是空时返回的错误
type ErrParamIsEmpty struct {
	name string
}

func (e *ErrParamIsEmpty) Error() string {
	return fmt.Sprintf("param %s shuld not be empty", e.name)
}

func NewErrParamIsEmpty(param string) error {
	return &ErrParamIsEmpty{name: param}
}

func IsErrParamIsEmpty(err error) bool {
	_, ok := err.(*ErrParamIsEmpty)
	return ok
}
