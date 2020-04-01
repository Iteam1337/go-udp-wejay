package utils

import "reflect"

func SetNil(i interface{}) {
	v := reflect.ValueOf(i)
	v.Elem().Set(reflect.Zero(v.Elem().Type()))
}
