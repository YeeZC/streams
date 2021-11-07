package slices

import (
	"fmt"
	"reflect"
)

func Contains(arr interface{}, elem interface{}) bool {
	switch arr := arr.(type) {
	case []string:
		return ContainsString(arr, elem.(string))
	case []int:
		return ContainsInt(arr, elem.(int))
	case []int64:
		return ContainsInt64(arr, elem.(int64))
	case []float64:
		return ContainsFloat64(arr, elem.(float64))
	default:
		return ContainsInf(arr, elem)
	}
}

func ContainsString(arr []string, elem string) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func ContainsInt(arr []int, elem int) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func ContainsInt64(arr []int64, elem int64) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func ContainsFloat64(arr []float64, elem float64) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

func ContainsInf(arr, elem interface{}) bool {
	arrValue := reflect.ValueOf(arr)
	if arrValue.Kind() != reflect.Array && arrValue.Kind() != reflect.Slice {
		panic(fmt.Errorf("%v is not array", arrValue.Type()))
	}
	for i := 0; i < arrValue.Len(); i++ {
		if reflect.DeepEqual(arrValue.Index(i).Interface(), elem) {
			return true
		}
	}
	return false
}
