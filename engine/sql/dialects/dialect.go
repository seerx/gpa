package dialects

import (
	"reflect"

	"github.com/seerx/gpa/engine/sql/types"
)

type Dialect interface {
	Init(*URI) error
	// 给字符串中的双引号转义
	QuoteExpr(str string) string
	// DataTypeOf 转为 sql 类型
	DataTypeOf(val reflect.Value) *types.SQLType
}
