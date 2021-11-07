package util

type T interface{}
type R interface{}

type Predicate func(T) bool
type Function func(T) R
type Comparator func(T, T) int
type BinaryOperator func(T, T) R
type Consumer func(T)
