package schema

import (
	"github.com/seerx/gpa/engine/objs"
	"github.com/seerx/gpa/engine/sql/names"
	"github.com/seerx/gpa/engine/sql/types"
)

type Column struct {
	// tag string

	fieldName string // 从 Field.Name 转换而来
	Field     objs.Object
	Type      types.SQLType

	// 忽略该字段 -
	Ignore bool
	// 主键 pk
	IsPrimaryKey bool
	// 自增类型 autoincr
	IsAutoIncrement bool
	// 是否可为空 null | not-null
	Nullable bool
	// 是否索引 index | index(索引名称)
	IsIndex bool
	// IndexNames []string
	// 是否唯一索引 unique | unique(索引名称)
	IsUnique bool
	// UniqueNames []string
	// 默认值 default(值)
	Default string
	// 添加时，写入当前时间
	// Created bool
	// 修改是，写入当前时间
	// Updated bool

	// IsJSON 是否 json
	IsJSON bool

	// EnumOptions map[string]int
	// SetOptions  map[string]int

	Length  int
	Length2 int

	Indexes map[string]int
}

func (c *Column) FieldName() string {
	if c.fieldName == "" {
		c.fieldName = names.ToTableFieldName(c.Field.Name)
	}
	return c.fieldName
}
