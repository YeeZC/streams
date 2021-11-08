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
	case []int32:
		return ContainsInt32(arr, elem.(int32))
	case []int64:
		return ContainsInt64(arr, elem.(int64))
	case []float32:
		return ContainsFloat32(arr, elem.(float32))
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

func ContainsInt32(arr []int32, elem int32) bool {
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

func ContainsFloat32(arr []float32, elem float32) bool {
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

func IndexOf(arr interface{}, elem interface{}) int {
	switch arr := arr.(type) {
	case []string:
		return IndexOfString(arr, elem.(string))
	case []int:
		return IndexOfInt(arr, elem.(int))
	case []int32:
		return IndexOfInt32(arr, elem.(int32))
	case []int64:
		return IndexOfInt64(arr, elem.(int64))
	case []float32:
		return IndexOfFloat32(arr, elem.(float32))
	case []float64:
		return IndexOfFloat64(arr, elem.(float64))
	default:
		return IndexOfInf(arr, elem)
	}
}

func IndexOfString(arr []string, elem string) int {
	for i, v := range arr {
		if v == elem {
			return i
		}
	}
	return -1
}

func IndexOfInt(arr []int, elem int) int {
	for i, v := range arr {
		if v == elem {
			return i
		}
	}
	return -1
}

func IndexOfInt32(arr []int32, elem int32) int {
	for i, v := range arr {
		if v == elem {
			return i
		}
	}
	return -1
}

func IndexOfInt64(arr []int64, elem int64) int {
	for i, v := range arr {
		if v == elem {
			return i
		}
	}
	return -1
}

func IndexOfFloat32(arr []float32, elem float32) int {
	for i, v := range arr {
		if v == elem {
			return i
		}
	}
	return -1
}

func IndexOfFloat64(arr []float64, elem float64) int {
	for i, v := range arr {
		if v == elem {
			return i
		}
	}
	return -1
}

func IndexOfInf(arr, elem interface{}) int {
	arrValue := reflect.ValueOf(arr)
	if arrValue.Kind() != reflect.Array && arrValue.Kind() != reflect.Slice {
		panic(fmt.Errorf("%v is not array", arrValue.Type()))
	}
	for i := 0; i < arrValue.Len(); i++ {
		if reflect.DeepEqual(arrValue.Index(i).Interface(), elem) {
			return i
		}
	}
	return -1
}
