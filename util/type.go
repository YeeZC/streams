package util

type Predicate func(interface{}) bool
type Function func(interface{}) interface{}
type Comparator func(interface{}, interface{}) int
type BinaryOperator func(interface{}, interface{}) interface{}
type Consumer func(interface{})
