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

type Type struct {
	Package string
	Name    string
	isPtr   bool
	isSlice bool
	// isMap   bool
	typ TypeClass
}

func newType(pkg, name string, isPtr bool, typ TypeClass) *Type {
	return &Type{Package: pkg, Name: name, isPtr: isPtr, typ: typ}
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
		return NewPrimitiveType(tp.Name)
	case *ast.StarExpr: // 指针类型
		raw := tp.X.(*ast.SelectorExpr)
		return parseType(raw, true)
	case *ast.SelectorExpr: // 非指针类型
		return parseType(tp, false)
		// TODO BLOB 类型待确定
	case *ast.ArrayType:
		typ := NewTypeByExpr(tp.Elt)
		if typ != nil {
			typ.isSlice = true
		}
		return typ
	default:
		fmt.Println(tp)
	}
	return nil
}

func parseType(expr *ast.SelectorExpr, ptr bool) *Type {
	x, ok := expr.X.(*ast.Ident)
	if ok {
		// 有 package ，认为是自定义类型
		if ptr {
			return NewPtrType(x.Name, expr.Sel.Name)
		}
		return NewType(x.Name, expr.Sel.Name)
		// return &Type{
		// 	Name:    expr.Sel.Name,
		// 	Package: x.Name,
		// 	isPtr:   ptr,
		// }
	}
	// 无 package ，认为是原生类型
	if ptr {
		return NewPtrPrimitiveType(expr.Sel.Name)
	}
	return NewPrimitiveType(expr.Sel.Name)
	// return NewType(x.Name, expr.Sel.Name)
	// return &Type{
	// 	Name:  expr.Sel.Name,
	// 	isPtr: ptr,
	// }
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
		isPtr:   ptr,
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

func NewPrimitiveType(name string) *Type {
	return newType("", name, false, PRIMITIVE)
}

func NewPtrPrimitiveType(name string) *Type {
	return newType("", name, true, PRIMITIVE)
}

func NewFuncType() *Type {
	return newType("", "func", false, FUNC)
}

func NewContextType() *Type {
	return newType("", "context", false, CONTEXT)
}

func NewTimeType() *Type {
	return newType("time", "Time", false, TIME)
}

func NewPtrTimeType() *Type {
	return newType("time", "Time", true, TIME)
}

func NewErrorType() *Type {
	return newType("errors", "Error", false, ERROR)
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
		typ.isPtr == t.isPtr
}

func (typ *Type) IsInt64() bool {
	return typ.Name == "int64" && typ.typ == PRIMITIVE
}

func (typ *Type) IsByte() bool {
	return typ.Name == "byte" && typ.typ == PRIMITIVE
}

func (typ *Type) IsTime() bool {
	return typ.typ == TIME
}

func (typ *Type) IsStruct() bool {
	return typ.Package != "" && typ.typ == CUSTOM
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

func (typ *Type) String() string {
	if typ.Package == "" {
		return typ.Name
	}
	return fmt.Sprintf("%s.%s", typ.Package, typ.Name)
}

func (typ *Type) StringExt() string {
	ptrTag := ""
	if typ.isPtr {
		ptrTag = "*"
	}
	if typ.Package == "" {
		return ptrTag + typ.Name
	}
	return fmt.Sprintf("%s%s.%s", ptrTag, typ.Package, typ.Name)
}
