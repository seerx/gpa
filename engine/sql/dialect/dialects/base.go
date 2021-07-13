package dialects

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/seerx/gpa/engine/sql/dialect/intf"
	"github.com/seerx/gpa/engine/sql/metas/schema"
	"github.com/seerx/gpa/engine/sql/types"
)

type baseDialect struct {
	intf.Dialect
	uri    intf.URI
	quoter intf.Quoter
}

func (bd *baseDialect) Init(dialect intf.Dialect, uri *intf.URI) error {
	bd.Dialect, bd.uri = dialect, *uri
	return nil
}
func (bd *baseDialect) URI() *intf.URI {
	return &bd.uri
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

// SQLDropTable 生成删除表的 SQL
func (bd *baseDialect) SQLDropTable(tableName string) (string, error) {
	quote := bd.Quoter().Quote
	tableName = bd.TableNameWithSchema(tableName)
	return fmt.Sprintf("DROP TABLE IF EXISTS %s", quote(tableName)), nil
}

func (bd *baseDialect) SQLColumn(col *schema.Column, inlinePrimaryKey bool) (string, error) {
	sql := strings.Builder{}
	// 字段名称
	if err := bd.Quoter().QuoteTo(&sql, col.FieldName()); err != nil {
		return "", err
	}
	if err := sql.WriteByte(' '); err != nil {
		return "", err
	}
	// 数据类型
	if _, err := sql.WriteString(bd.SQLType(col)); err != nil {
		return "", err
	}
	if err := sql.WriteByte(' '); err != nil {
		return "", err
	}

	if inlinePrimaryKey && col.IsPrimaryKey {
		// 只有一个字段是主键，且该字段是主键
		if _, err := sql.WriteString("PRIMARY KEY "); err != nil {
			return "", err
		}

		if col.IsAutoIncrement {
			// 该字段是自增类型
			if _, err := sql.WriteString(bd.AutoIncrStr()); err != nil {
				return "", err
			}
			if err := sql.WriteByte(' '); err != nil {
				return "", err
			}
		}
	}

	if col.Default != "" {
		if _, err := sql.WriteString("DEFAULT "); err != nil {
			return "", err
		}
		if _, err := sql.WriteString(col.Default); err != nil {
			return "", err
		}
		if err := sql.WriteByte(' '); err != nil {
			return "", err
		}
	}

	if col.Nullable {
		if _, err := sql.WriteString("NULL "); err != nil {
			return "", err
		}
	} else {
		if _, err := sql.WriteString("NOT NULL "); err != nil {
			return "", err
		}
	}

	return sql.String(), nil
}

func (bd *baseDialect) SQLAddColumn(tableName string, col *schema.Column) string {
	s, _ := bd.SQLColumn(col, true)
	tableName = bd.TableNameWithSchema(tableName)
	return fmt.Sprintf("ALTER TABLE %v ADD %v", bd.Quoter().Quote(tableName), s)
}

func (bd *baseDialect) SQLModifyColumn(tableName string, col *schema.Column) string {
	s, _ := bd.SQLColumn(col, false)
	return fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s", tableName, s)
}

func (bd *baseDialect) SQLCreateIndex(tableName string, index *schema.Index) string {
	quoter := bd.Dialect.Quoter()
	var unique string
	var idxName string
	if index.Type == types.UniqueType {
		unique = " UNIQUE"
	}
	idxName = index.XName(tableName)
	return fmt.Sprintf("CREATE%s INDEX %v ON %v (%v)", unique,
		quoter.Quote(idxName), quoter.Quote(tableName),
		quoter.Join(index.Cols, ","))
}

func (bd *baseDialect) SQLDropIndex(tableName string, index *schema.Index) string {
	quote := bd.Dialect.Quoter().Quote
	var name string
	if index.Regular {
		name = index.XName(tableName)
	} else {
		name = index.Name
	}
	return fmt.Sprintf("DROP INDEX %v ON %s", quote(name), quote(tableName))
}
