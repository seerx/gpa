package xtype

import (
	"fmt"
	"time"

	"github.com/seerx/gpa/engine/sql/metas/schema"
)

type Field struct {
	schema.Column

	VarName string // 对应的 golang 中的变量名称
	ArgName string // 对应的参数名称
	// SQLType *types.SQLType
	TimeZone *time.Location
}

type ParamType struct {
	File      string   // 所在文件名称
	Name      string   // 名称
	TableName string   // 对应的数据库表名
	Fields    []*Field // 结构体对应的成员列表
}

var paramsPool = map[string]map[string]*ParamType{}

func GetParam(name, dir string) (*ParamType, error) {
	var err error
	params, ok := paramsPool[dir]
	if !ok {
		params, err = Scan(dir)
		if err != nil {
			return nil, err
		}
		paramsPool[dir] = params
	}

	param, ok := params[name]
	if ok {
		return param, nil
	}
	return nil, fmt.Errorf("no struct %s is defined in %s", name, dir)
}
