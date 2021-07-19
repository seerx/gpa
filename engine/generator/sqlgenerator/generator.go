package sqlgenerator

import "github.com/seerx/gpa/rt/dbutil"

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
	TableName                string
	Columns                  []string
	SelectFields             []string
	ParamPlaceHolder         []string
	Params                   []*SQLParam
	Where                    string
	WhereParams              []*SQLParam
	ReturnAutoincrPrimaryKey string // insert 时是否返回自增主键
}

type SQLGenerator interface {
	Insert(sql *SQL) (string, []*SQLParam)
	Update(sql *SQL) (string, []*SQLParam, []*SQLParam)
	Delete(sql *SQL) (string, []*SQLParam)
	Query(sql *SQL) (string, []*SQLParam)
}
