package utils

import "reflect"

func ReflectValueAndType(bean interface{}) (reflect.Value, reflect.Type) {
	val := reflect.Indirect(reflect.ValueOf(bean))
	return val, val.Type()
}
