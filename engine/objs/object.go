package objs

import (
	"reflect"

	"github.com/seerx/gpa/engine/sql/types"
)

type Object struct {
	Name  string // 名称
	Type  Type   // 类型
	Key   *Type  // map 时 key 的类型
	slice bool   // 是否数组

	Params    []*Object
	ParamsMap map[string]*Object
	Results   []*Object
}

func NewObjectFromStructField(field *reflect.StructField) *Object {
	typ := NewTypeFromStructField(field)
	return &Object{
		Name:  field.Name,
		Type:  *typ,
		slice: typ.isArray,
	}
}

func NewObject(name string, typ Type) *Object { return &Object{Name: name, Type: typ} }
func NewSliceObject(name string, typ Type) *Object {
	return &Object{Name: name, Type: typ, slice: true}
}
func NewMapObject(name string, typ, key Type) *Object {
	return &Object{Name: name, Type: typ, Key: &key}
}

func (o *Object) IsSlice() bool { return o.slice }
func (o *Object) IsMap() bool   { return o.Key != nil }

func (o *Object) AddParam(obj *Object) {
	obj.Params = append(obj.Params, obj)
	if obj.ParamsMap == nil {
		obj.ParamsMap = map[string]*Object{}
	}
	obj.ParamsMap[obj.Name] = obj
}

func (o *Object) AddResult(obj *Object) { o.Results = append(o.Results, obj) }

func (o *Object) GetSQLType() (st *types.SQLType) {
	if o.IsMap() || o.IsSlice() {
		// map 和数组
		if o.Type.IsByte() {
			st = &types.SQLType{types.Blob, 0, 0}
		} else {
			st = &types.SQLType{types.Text, 0, 0}
		}
	} else if o.Type.IsPrimitive() {
		// 基础类型
		switch o.Type.Name {
		case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32":
			st = &types.SQLType{types.Int, 0, 0}
		case "int64", "uint64":
			st = &types.SQLType{types.BigInt, 0, 0}
		case "float32":
			st = &types.SQLType{types.Float, 0, 0}
		case "float64":
			st = &types.SQLType{types.Double, 0, 0}
		case "complex64", "complex128":
			st = &types.SQLType{types.Varchar, 64, 0}
		case "bool":
			st = &types.SQLType{types.Bool, 0, 0}
		case "string":
			st = &types.SQLType{types.Varchar, 255, 0}
		}
	} else if o.Type.IsTime() {
		// 时间
		st = &types.SQLType{types.DateTime, 0, 0}
	} else if o.Type.IsStruct() {
		// 自定义结构
		st = &types.SQLType{types.Text, 0, 0}
	} else {
		// 未知类型
		st = &types.SQLType{types.Text, 0, 0}
	}
	return
}
