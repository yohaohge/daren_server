package store

import (
	"reflect"
)

// obj should be a object not a pointer
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		var val = obj2.Field(i).Interface()
		//if obj2.Field(i).Type() == reflect.TypeOf(FlexInt(0)) {
		//	val = obj2.Field(i).Int()
		//}
		data[obj1.Field(i).Tag.Get("json")] = val
	}
	return data
}
