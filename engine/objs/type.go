package objs

import (
	"context"
	"fmt"
	"go/ast"
	"reflect"
	"time"

	"github.com/seerx/gpa/rt/dbutil"
)

type TypeClass int

const (
	PRIMITIVE TypeClass = iota // 基础类型
	FUNC                       // 函数
	TIME                       // 时间
	CONTEXT                    // context.Context
	ERROR                      // error
	BLOB                       // BLOB 读写类型
	CUSTOM                     // 自定义类型
)

var primitiveTypes = map[string]bool{
	"int":    true,
	"int8":   true,
	"int16":  true,
	"int32":  true,
	"int64":  true,
	"uint":   true,
	"uint8":  true,
	"uint16": true,
	"uint32": true,
	"uint64": true,

	"byte":   true,
	"bool":   true,
	"string": true,

	"float32": true,
	"float64": true,
}

type Type struct {
	Package string
	Name    string
	IsPtr   bool
	isSlice bool
	// isMap   bool
	typ TypeClass
}

func newType(pkg, name string, isPtr bool, typ TypeClass) *Type {
	return &Type{Package: pkg, Name: name, IsPtr: isPtr, typ: typ}
}

var (
	_time         time.Time
	typeOfTime    = reflect.TypeOf(_time)
	typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()
	typeOfError   = reflect.TypeOf((*error)(nil)).Elem()
	typeOfBlob    = reflect.TypeOf((*dbutil.BlobReadWriter)(nil)).Elem()
)

func NewTypeByExpr(expr ast.Expr) *Type {
	switch tp := expr.(type) {
	case *ast.Ident: // 原生类型
		return NewTypeByPkgAndName("", tp.Name)
	case *ast.StarExpr: // 指针类型
		raw, ok := tp.X.(*ast.SelectorExpr)
		if ok {
			return parseType(raw, true)
		} else {
			tmp, ok := tp.X.(ast.Expr)
			if ok {
				tp := tmp.(*ast.Ident)
				return NewPtrTypeByPkgAndName("", tp.Name)
			}
		}
	case *ast.SelectorExpr: // 非指针类型
		return parseType(tp, false)
	case *ast.ArrayType:
		typ := NewTypeByExpr(tp.Elt)
		if typ != nil {
			typ.isSlice = true
		}
		return typ
	default:
		// fmt.Println(tp)
	}
	return nil
}

func parseType(expr *ast.SelectorExpr, ptr bool) *Type {
	x, ok := expr.X.(*ast.Ident)
	pkg := ""
	if ok {
		// 没有 package
		pkg = x.Name
	}
	// if ok {
	if ptr {
		return NewPtrTypeByPkgAndName(pkg, expr.Sel.Name)
	}
	return NewTypeByPkgAndName(pkg, expr.Sel.Name)
}

func NewTypeFromStructField(field *reflect.StructField) *Type {
	typ := field.Type
	// kind := typ.Kind()
	slice := false
	ptr := false
	if typ.Kind() == reflect.Slice {
		// 数组
		slice = true
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Ptr {
		// 是指针
		ptr = true
		typ = typ.Elem()
	}

	res := &Type{
		Package: typ.PkgPath(),
		Name:    typ.Name(),
		IsPtr:   ptr,
		isSlice: slice,
	}

	// typ.ConvertibleTo()
	if typ == typeOfTime {
		res.typ = TIME
	} else if typ.ConvertibleTo(typeOfContext) {
		res.typ = CONTEXT
	} else if typ.ConvertibleTo(typeOfError) {
		res.typ = ERROR
	} else if typ.ConvertibleTo(typeOfBlob) {
		res.typ = BLOB
	} else if res.Package == "" {
		res.typ = PRIMITIVE
	} else {
		res.typ = CUSTOM
	}

	return res
}

func NewType(pkg, name string) *Type {
	return newType(pkg, name, false, CUSTOM)
}

func NewPtrType(pkg, name string) *Type {
	return newType(pkg, name, true, CUSTOM)
}

// func NewPrimitiveType(name string) (*Type, bool) {
// 	if _, ok := primitiveTypes[name]; ok {
// 		return newType("", name, false, PRIMITIVE), true
// 	}
// 	return nil, false
// }

// func NewPtrPrimitiveType(name string) (*Type, bool) {
// 	if _, ok := primitiveTypes[name]; ok {
// 		return newType("", name, true, PRIMITIVE), true
// 	}
// 	return nil, false
// }

func NewFuncType() *Type {
	return newType("", "func", false, FUNC)
}

func NewTypeByPkgAndName(pkg, name string) *Type {
	if pkg == "context" && name == "Context" {
		return NewContextType()
	}
	if pkg == "time" && name == "Time" {
		return NewTimeType()
	}
	if pkg == "" {
		if name == "error" {
			return NewErrorType()
		}
		if _, ok := primitiveTypes[name]; ok {
			// 原生类型
			return newType("", name, false, PRIMITIVE)
		}
	}

	return newType(pkg, name, false, CUSTOM)
}

func NewPtrTypeByPkgAndName(pkg, name string) *Type {
	if pkg == "context" && name == "Context" {
		return NewContextType()
	}
	if pkg == "time" && name == "Time" {
		return NewPtrTimeType()
	}
	if pkg == "" {
		if _, ok := primitiveTypes[name]; ok {
			// 原生类型
			return newType("", name, true, PRIMITIVE)
		}
	}

	return newType(pkg, name, true, CUSTOM)
}

func NewContextType() *Type {
	return newType("context", "Context", false, CONTEXT)
}

func NewTimeType() *Type {
	return newType("time", "Time", false, TIME)
}

func NewPtrTimeType() *Type {
	return newType("time", "Time", true, TIME)
}

func NewErrorType() *Type {
	return newType("", "error", false, ERROR)
}

// Equals 判断是否相同类型，不区分是否指针
func (typ *Type) Equals(t *Type) bool {
	return typ.typ == t.typ &&
		typ.Package == t.Package &&
		typ.Name == t.Name
}

// Equals 判断是否相同类型，包含是否指针
func (typ *Type) EqualsExactly(t *Type) bool {
	return typ.typ == t.typ &&
		typ.Package == t.Package &&
		typ.Name == t.Name &&
		typ.IsPtr == t.IsPtr
}

func (typ *Type) IsInt64() bool {
	return typ.Name == "int64" && typ.typ == PRIMITIVE
}

func (typ *Type) IsGenericInt() bool {
	return typ.typ == PRIMITIVE &&
		(typ.Name == "int" ||
			typ.Name == "uint" ||
			typ.Name == "int64" ||
			typ.Name == "int32" ||
			typ.Name == "int16" ||
			typ.Name == "int8" ||
			typ.Name == "uint64" ||
			typ.Name == "uint32" ||
			typ.Name == "uint16" ||
			typ.Name == "uint8")
}

func (typ *Type) IsByte() bool {
	return typ.Name == "byte" && typ.typ == PRIMITIVE
}

func (typ *Type) IsTime() bool {
	return typ.typ == TIME
}

func (typ *Type) IsContext() bool {
	return typ.typ == CONTEXT
}

func (typ *Type) IsError() bool {
	return typ.typ == ERROR
}

func (typ *Type) IsPrimitive() bool {
	return typ.typ == PRIMITIVE
}

func (typ *Type) IsCustom() bool {
	return typ.typ == CUSTOM
}

func (typ *Type) String() string {
	if typ.Package == "" {
		return typ.Name
	}
	return fmt.Sprintf("%s.%s", typ.Package, typ.Name)
}

func (typ *Type) StringExt() string {
	ptrTag := ""
	if typ.IsPtr {
		ptrTag = "*"
	}
	if typ.Package == "" {
		return ptrTag + typ.Name
	}
	return fmt.Sprintf("%s%s.%s", ptrTag, typ.Package, typ.Name)
}
