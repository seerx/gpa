package objs

import (
	"errors"
	"fmt"
	"go/ast"
	"reflect"
	"strings"

	"github.com/seerx/gpa/engine/sql/types"
)

type Object struct {
	Name      string // 名称
	Type      Type   // 类型
	Key       *Type  // map 时 key 的类型
	IsSlice   bool   // 是否数组
	IsMap     bool
	IsFunc    bool
	Params    []*Object
	Results   []*Object
	ParamsMap map[string]*Object
	Extra     interface{}
}

func NewObjectFromStructField(field *reflect.StructField) *Object {
	typ := NewTypeFromStructField(field)
	return &Object{
		Name:      field.Name,
		Type:      *typ,
		IsSlice:   typ.isSlice,
		ParamsMap: map[string]*Object{},
	}
}

func NewObject(name string, typ Type) *Object {
	return &Object{
		Name:      name,
		Type:      typ,
		IsSlice:   typ.isSlice,
		ParamsMap: map[string]*Object{},
	}
}

func NewEmptyObject() *Object {
	return &Object{
		ParamsMap: map[string]*Object{},
	}
}

// func NewSliceObject(name string, typ Type) *Object {
// 	return &Object{Name: name, Type: typ, slice: true}
// }
// func NewMapObject(name string, typ, key Type) *Object {
// 	return &Object{Name: name, Type: typ, Key: &key}
// }

// func (o *Object) IsSlice() bool { return o.slice }
// func (o *Object) IsMap() bool   { return o.Key != nil }

func (o *Object) FindParam(names []string) *Object {
	for _, name := range names {
		p, ok := o.ParamsMap[name]
		if ok {
			return p
		}
	}
	return nil
}

func (o *Object) AddParam(obj *Object) {
	o.Params = append(o.Params, obj)
	if o.ParamsMap == nil {
		o.ParamsMap = map[string]*Object{}
	}
	o.ParamsMap[obj.Name] = obj
}

func (o *Object) AddResult(obj *Object) { o.Results = append(o.Results, obj) }

func (o *Object) GetSQLTypeByType() (st *types.SQLType) {
	if o.IsMap || o.IsSlice {
		// map 和数组
		if o.Type.IsByte() {
			st = &types.SQLType{Name: types.Blob, Length: 0, Length2: 0}
		} else {
			st = &types.SQLType{Name: types.Text, Length: 0, Length2: 0}
		}
	} else if o.Type.IsPrimitive() {
		// 基础类型
		switch o.Type.Name {
		case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32":
			st = &types.SQLType{Name: types.Int, Length: 0, Length2: 0}
		case "int64", "uint64":
			st = &types.SQLType{Name: types.BigInt, Length: 0, Length2: 0}
		case "float32":
			st = &types.SQLType{Name: types.Float, Length: 0, Length2: 0}
		case "float64":
			st = &types.SQLType{Name: types.Double, Length: 0, Length2: 0}
		case "complex64", "complex128":
			st = &types.SQLType{Name: types.Varchar, Length: 64, Length2: 0}
		case "bool":
			st = &types.SQLType{Name: types.Bool, Length: 0, Length2: 0}
		case "string":
			st = &types.SQLType{Name: types.Varchar, Length: 255, Length2: 0}
		}
	} else if o.Type.IsTime() {
		// 时间
		st = &types.SQLType{Name: types.DateTime, Length: 0, Length2: 0}
	} else if o.Type.IsCustom() {
		// 自定义结构
		st = &types.SQLType{Name: types.Text, Length: 0, Length2: 0}
	} else {
		// 未知类型
		st = &types.SQLType{Name: types.Text, Length: 0, Length2: 0}
	}
	return
}

func (o *Object) ParseFunc(params *ast.FieldList,
	results *ast.FieldList,
	dialect string,
	fnFuncParseCb func(*Object) error) error {
	// o.Name = GetName(method.Names) // 函数名称
	// // f.SQL = ParseSQL(method.Doc.Text(), dialect) // SQL 语句定义

	// typ, ok := method.Type.(*ast.FuncType)
	// if !ok {
	// 	return fmt.Errorf("%s is not a valid method", o.Name)
	// }

	// 遍历参数列表
	if params != nil {
		for _, mp := range params.List {
			param := NewEmptyObject()
			param.Name = GetName(mp.Names)
			if err := param.Parse(mp, mp.Type, dialect, fnFuncParseCb, 0); err != nil {
				return err
			}
			o.AddParam(param)
		}
	}
	// 遍历返回值列表
	if results != nil {
		for _, p := range results.List {
			result := NewEmptyObject() // NewObject(f.repo) // {Name: getName(p.Names)}
			result.Name = GetName(p.Names)
			if err := result.Parse(p, p.Type, dialect, fnFuncParseCb, 0); err != nil {
				return err
			}
			o.AddResult(result)
		}
	}

	return nil
}

func (o *Object) Parse(field *ast.Field,
	expr ast.Expr,
	dialect string,
	fnFuncParseCb func(*Object) error,
	level int) error {
	var err error
	switch pt := expr.(type) {
	case *ast.Ident:
		// 普通类型
		// NewTypeByPkgAndName()
		o.Type = *NewTypeByPkgAndName("", pt.Name)
		// if pt.Name == "error" {
		// 	o.Type = *NewErrorType()
		// } else {
		// 	o.Type = *NewPrimitiveType(pt.Name)
		// }
	case *ast.FuncType:
		// 函数类型
		// if err := fnFuncParseCb(o, level); err != nil {
		// 	return err
		// }
		o.Name = GetName(field.Names) // 函数名称
		typ, ok := field.Type.(*ast.FuncType)
		if !ok {
			return fmt.Errorf("%s is not a valid method", o.Name)
		}
		o.IsFunc = true
		if err := o.ParseFunc(typ.Params, typ.Results, dialect, fnFuncParseCb); err != nil {
			return err
		}
		if fnFuncParseCb != nil {
			if err := fnFuncParseCb(o); err != nil {
				return err
			}
		}

		// f.SQL = ParseSQL(method.Doc.Text(), dialect) // SQL 语句定义

		// if level > 1 {
		// 	return errors.New("不支持多层嵌套函数类型")
		// }

		// fn := NewFuncWithObject(o)
		// if err := fn.Parse(field, dialect); err != nil {
		// 	return err
		// }
		// o.Type = *objs.NewFuncType()
	case *ast.SelectorExpr:
		o.Type = *ParseSelectorExprType(pt, false)
	case *ast.StarExpr:
		o.Type = *ParseSelectorExprType(pt.X.(*ast.SelectorExpr), true)
	case *ast.ArrayType:
		// if err := fnSliceParseCb(o, pt.Elt, level); err != nil {
		// 	return err
		// }
		o.IsSlice = true
		if level > 1 {
			return errors.New("不支持多层嵌套数据类型")
		}
		obj := NewEmptyObject()
		if err := obj.Parse(nil, pt.Elt, dialect, fnFuncParseCb, level+1); err != nil {
			return err
		}
		o.Type = obj.Type
		// obj := NewObject(o.repo)
		// if err := obj.Parse(nil, pt.Elt, dialect, level+1); err != nil {
		// 	return err
		// }
		// o.Object = obj.Object
		// o.IsSlice = true
		// arg.Slice = true
	case *ast.SliceExpr:
		// if err := fnSliceParseCb(o, pt.X, level); err != nil {
		// 	return err
		// }
		o.IsSlice = true
		obj := NewEmptyObject()
		if err := obj.Parse(nil, pt.X, dialect, fnFuncParseCb, level+1); err != nil {
			return err
		}
		o.Type = obj.Type
		// if level > 1 {
		// 	return errors.New("不支持多层嵌套数据类型")
		// }
		// obj := NewObject(o.repo)
		// if err := obj.Parse(nil, pt.X, dialect, level+1); err != nil {
		// 	return err
		// }
		// o.Object = obj.Object
		// o.IsSlice = true
	case *ast.MapType:
		// if err := fnMapParseCb(o, pt.Key, pt.Value, level); err != nil {
		// 	return err
		// }
		// if level > 1 {
		// 	return errors.New("不支持多层嵌套数据类型")
		// }
		o.IsMap = true
		obj := NewEmptyObject()
		if err := obj.Parse(nil, pt.Value, dialect, fnFuncParseCb, level+1); err != nil {
			return err
		}
		o.Type = obj.Type
		// o.Object = obj.Object
		// obj := NewObject(o.repo)
		// if err := obj.Parse(nil, pt.Value, dialect, level+1); err != nil {
		// 	return err
		// }
		// o.Object = obj.Object

		obj = NewEmptyObject()
		if err := obj.Parse(nil, pt.Key, dialect, fnFuncParseCb, level+1); err != nil {
			return err
		}
		o.Key = &obj.Type
		// obj = NewObject(o.repo)
		// if err := obj.Parse(nil, pt.Key, dialect, level+1); err != nil {
		// 	return err
		// }
		// o.Key = &obj.Type
	default:
		err = errors.New("不支持的数据类型")
	}
	return err
}

func ParseSelectorExprType(se *ast.SelectorExpr, ptr bool) *Type {
	// name := ""
	x, ok := se.X.(*ast.Ident)
	pkg := ""
	if ok {
		pkg = x.Name
	}
	if ptr {
		return NewPtrTypeByPkgAndName(pkg, se.Sel.Name)
	}
	return NewTypeByPkgAndName(pkg, se.Sel.Name)
}

func GetName(names []*ast.Ident) string {
	for _, name := range names {
		return name.Name
	}
	return ""
}

func ParseSQL(comment, dialect string) string {
	lines := strings.Split(comment, "\n")
	sql := ""
	dialect += ":"
	for _, line := range lines {
		if strings.Index(line, dialect) == 0 {
			sql = line[len(dialect):]
			break
		}
		if strings.Index(line, "sql:") == 0 {
			sql = line[len("sql:"):]
		}
	}
	return sql
}
