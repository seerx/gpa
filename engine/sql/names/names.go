package names

import "reflect"

// GetTableName 对象转表名称
func GetTableName(v reflect.Value) string {
	return getTableName(LintGonicMapper, v)
}

// ToTableName 对象名称转表名，成员名称转字段名称
func ToTableName(name string) string {
	return LintGonicMapper.Obj2Table(name)
}

// ToObjName 对象名称转字段名称
func ToObjName(name string) string {
	return LintGonicMapper.Table2Obj(name)
}
