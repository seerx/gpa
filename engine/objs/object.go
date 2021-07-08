package objs

import "reflect"

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
