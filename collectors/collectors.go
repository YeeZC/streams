package collectors

import (
	"reflect"

	"github.com/yeezc/streams/util"
)

type Collector interface {
	Supplier() util.Supplier
	Combiner() util.BinaryOperator
	Finisher() util.Function
}

type collectorImpl struct {
	supplier util.Supplier
	combiner util.BinaryOperator
	finisher util.Function
}

func (c *collectorImpl) Supplier() util.Supplier {
	return c.supplier
}
func (c *collectorImpl) Combiner() util.BinaryOperator {
	return c.combiner
}
func (c *collectorImpl) Finisher() util.Function {
	return c.finisher
}

func Of(s util.Supplier, c util.BinaryOperator, finisher util.Function) Collector {
	return &collectorImpl{s, c, finisher}
}

func ToSlice(t reflect.Type) Collector {
	return Of(func() interface{} {
		return reflect.MakeSlice(reflect.SliceOf(t), 0, 0)
	}, func(i1, i2 interface{}) interface{} {
		return reflect.Append(i1.(reflect.Value), reflect.ValueOf(i2))
	}, func(i interface{}) interface{} {
		return i.(reflect.Value).Interface()
	})
}

func ToMap(k, v reflect.Type, key, value util.Function) Collector {
	return Of(func() interface{} {
		return reflect.MakeMap(reflect.MapOf(k, v))
	}, func(i1, i2 interface{}) interface{} {
		mapKey, mapValue := reflect.ValueOf(key(i2)), reflect.ValueOf(value(i2))
		i1.(reflect.Value).SetMapIndex(mapKey, mapValue)
		return i1
	}, func(i interface{}) interface{} {
		return i.(reflect.Value).Interface()
	})
}
