package intf

import (
	"context"
	"reflect"

	"github.com/seerx/gpa/engine/sql/metas/schema"
	"github.com/seerx/gpa/engine/sql/types"
	"github.com/seerx/gpa/rt/dbutil"
	"github.com/seerx/gpa/rt/exec"
)

type SQLParam struct {
	VarName           string
	JSON              bool
	Time              bool
	Blob              bool // 实现了 BlobReadWriter
	SQLType           string
	VarAlias          string
	TimeProp          *dbutil.TimePropDesc
	SQLParamName      string // 从 SQL 语句中分析出来的参数名称，不作任何改变
	SQLParamFieldName string

	IsInOperator       bool
	InParamPlaceHolder string // in 操作符占位字符串
}

type SQL struct {
	TableName        string
	Columns          []string
	SelectFields     []string
	ParamPlaceHolder []string
	Params           []*SQLParam
	Where            string
	WhereParams      []*SQLParam
}

type Dialect interface {
	Init(*URI) error
	URI() *URI

	// 给字符串中的双引号转义
	QuoteExpr(str string) string
	// ToSQLType 转为 sql 类型
	ToSQLType(typ reflect.Type) *types.SQLType
	// DataTypeOf 转为 sql 类型
	DataTypeOf(val reflect.Value) *types.SQLType
	// SQLType 转换为 SQL 数据类型
	SQLType(col *schema.Column) string
	// AutoIncrStr 自增字段标志字符串
	AutoIncrStr() string

	Quoter() Quoter

	TableNameWithSchema(tableName string) string

	// 从数据库中查询所有表结构
	GetTables(ex exec.SQLExecutor, ctx context.Context) ([]*schema.Table, error)
	// SQLTableExists 判断表是否存在SQL
	SQLTableExists(tableName string) (sql string, args []interface{})
	// SQLCreateTable 生成创建表结构的 SQL
	SQLCreateTable(table *schema.Table, tableName string) ([]string, error)
	// SQLDropTable 生成删除表的 SQL
	SQLDropTable(tableName string) (string, error)

	GetColumns(ex exec.SQLExecutor, ctx context.Context, tableName string) ([]string, map[string]*schema.Column, error)
	// SQLColumn 生成列相关的 SQL 定义
	SQLColumn(col *schema.Column, inlinePrimaryKey bool) (string, error)
	SQLAddColumn(tableName string, col *schema.Column) string
	SQLModifyColumn(tableName string, col *schema.Column) string

	GetIndexes(ex exec.SQLExecutor, ctx context.Context, tableName string) (map[string]*schema.Index, error)
	SQLCreateIndex(tableName string, index *schema.Index) string
	SQLDropIndex(tableName string, index *schema.Index) string

	CreateInsertSQL(sql *SQL) (string, []*SQLParam)
	CreateUpdateSQL(sql *SQL) (string, []*SQLParam, []*SQLParam)
	CreateDeleteSQL(sql *SQL) (string, []*SQLParam)
	CreateQuerySQL(sql *SQL) (string, []*SQLParam)
}
