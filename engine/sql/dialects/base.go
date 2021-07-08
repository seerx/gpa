package dialects

import (
	"reflect"
	"strings"

	"github.com/seerx/gpa/engine/sql/types"
)

type baseDialect struct {
	Dialect
	uri URI
	// quoter Quoter
}

func (bd *baseDialect) Init(dialect Dialect, uri *URI) error {
	bd.Dialect, bd.uri = dialect, *uri
	return nil
}

func (bd *baseDialect) QuoteExpr(str string) string {
	return strings.ReplaceAll(str, "\"", "\\\"")
}

// ToSQLType 转为 sql 类型
func (bd *baseDialect) ToSQLType(typ reflect.Type) *types.SQLType {
	return types.Type2SQLType(typ)
}

func (bd *baseDialect) DataTypeOf(val reflect.Value) *types.SQLType {
	return types.Type2SQLType(val.Type())
}
