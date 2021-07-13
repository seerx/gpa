package types

import (
	"reflect"
	"time"
)

type SQLType struct {
	Name    string
	Length  int
	Length2 int
}

var (
	c_BYTE_DEFAULT byte
	c_TIME_DEFAULT time.Time
)

func Type2SQLType(t reflect.Type) (st *SQLType) {
	switch k := t.Kind(); k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		st = &SQLType{Int, 0, 0}
	case reflect.Int64, reflect.Uint64:
		st = &SQLType{BigInt, 0, 0}
	case reflect.Float32:
		st = &SQLType{Float, 0, 0}
	case reflect.Float64:
		st = &SQLType{Double, 0, 0}
	case reflect.Complex64, reflect.Complex128:
		st = &SQLType{Varchar, 64, 0}
	case reflect.Array, reflect.Slice, reflect.Map:
		if t.Elem() == reflect.TypeOf(c_BYTE_DEFAULT) {
			st = &SQLType{Blob, 0, 0}
		} else {
			st = &SQLType{Text, 0, 0}
		}
	case reflect.Bool:
		st = &SQLType{Bool, 0, 0}
	case reflect.String:
		st = &SQLType{Varchar, 255, 0}
	case reflect.Struct:
		if t.ConvertibleTo(reflect.TypeOf(c_TIME_DEFAULT)) {
			st = &SQLType{DateTime, 0, 0}
		} else {
			// TODO need to handle association struct
			st = &SQLType{Text, 0, 0}
		}
	case reflect.Ptr:
		st = Type2SQLType(t.Elem())
	default:
		st = &SQLType{Text, 0, 0}
	}
	return
}

func (s *SQLType) IsType(st int) bool {
	if t, ok := SqlTypes[s.Name]; ok && t == st {
		return true
	}
	return false
}

func (s *SQLType) IsText() bool {
	return s.IsType(TEXT_TYPE)
}

func (s *SQLType) IsBlob() bool {
	return s.IsType(BLOB_TYPE)
}

func (s *SQLType) IsTime() bool {
	return s.IsType(TIME_TYPE)
}

// Object2SQLType generate SQLType acorrding Go's type name(string)
// func Object2SQLType(obj *objs.Object) (st *SQLType) {
// 	if obj.IsMap() || obj.IsSlice() {
// 		// map 和数组
// 		if obj.Type.IsByte() {
// 			st = &SQLType{Blob, 0, 0}
// 		} else {
// 			st = &SQLType{Text, 0, 0}
// 		}
// 	} else if obj.Type.IsPrimitive() {
// 		// 基础类型
// 		switch obj.Type.Name {
// 		case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32":
// 			st = &SQLType{Int, 0, 0}
// 		case "int64", "uint64":
// 			st = &SQLType{BigInt, 0, 0}
// 		case "float32":
// 			st = &SQLType{Float, 0, 0}
// 		case "float64":
// 			st = &SQLType{Double, 0, 0}
// 		case "complex64", "complex128":
// 			st = &SQLType{Varchar, 64, 0}
// 		case "bool":
// 			st = &SQLType{Bool, 0, 0}
// 		case "string":
// 			st = &SQLType{Varchar, 255, 0}
// 		}
// 	} else if obj.Type.IsTime() {
// 		// 时间
// 		st = &SQLType{DateTime, 0, 0}
// 	} else if obj.Type.IsStruct() {
// 		// 自定义结构
// 		st = &SQLType{Text, 0, 0}
// 	} else {
// 		// 未知类型
// 		st = &SQLType{Text, 0, 0}
// 	}
// 	return
// }
