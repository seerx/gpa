package names

import "reflect"

// ToTableName 对象转表名称
func ToTableName(v reflect.Value) string {
	return GetTableName(LintGonicMapper, v)
}

// ToTableFieldName 对象名称转字段名称
func ToTableFieldName(name string) string {
	return LintGonicMapper.Obj2Table(name)
}

// ToObjName 对象名称转字段名称
func ToObjName(name string) string {
	return LintGonicMapper.Table2Obj(name)
}
