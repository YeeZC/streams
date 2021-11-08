package maps

import (
	"fmt"
	"reflect"
	"sync"
)

func Keys(i interface{}) []interface{} {
	if _map, ok := i.(sync.Map); ok {
		elems := make([]interface{}, 0)
		_map.Range(func(key, value interface{}) bool {
			elems = append(elems, key)
			return true
		})
		return elems
	}
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Map {
		panic(fmt.Errorf("%v is not a map", v.Type()))
	}
	keys := v.MapKeys()
	ret := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		ret = append(ret, key.Interface())
	}
	return ret
}

func Values(i interface{}) []interface{} {
	if _map, ok := i.(sync.Map); ok {
		elems := make([]interface{}, 0)
		_map.Range(func(key, value interface{}) bool {
			elems = append(elems, value)
			return true
		})
		return elems
	}
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Map {
		panic(fmt.Errorf("%v is not a map", v.Type()))
	}
	it := v.MapRange()
	elems := make([]interface{}, 0, v.Len())
	for it.Next() {
		elems = append(elems, it.Value().Interface())
	}
	return elems
}
