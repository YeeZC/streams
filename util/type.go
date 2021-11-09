package util

type Predicate func(interface{}) bool
type Function func(interface{}) interface{}
type BreakableFunction func(interface{}) (interface{}, bool)
type Comparator func(interface{}, interface{}) int
type BinaryOperator func(interface{}, interface{}) interface{}
type Consumer func(interface{})
type BreakableConsumer func(interface{}) bool
