package streams

import (
	"reflect"

	"github.com/yeezc/streams/util"
)

type T util.T
type R util.R

type Predicate util.Predicate
type Function util.Function
type Comparator util.Comparator
type BinaryOperator util.BinaryOperator
type Consumer util.Consumer

type Stream interface {
	Filter(predicate Predicate) Stream
	Map(function Function) Stream
	FindAny() util.Optional
	Distinct() Stream
	Sorted(c Comparator) Stream
	ForEach(consumer Consumer)
	Reduce(identity T, op BinaryOperator) R
	ToArray() interface{}
}

func Empty() Stream {
	c := make(chan interface{})
	close(c)
	return &defaultStream{in: c}
}

func Of(i interface{}) Stream {
	in := make(chan interface{}, 1)
	go func() {
		defer close(in)
		v := reflect.ValueOf(i)
		if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
			for i := 0; i < v.Len(); i++ {
				in <- v.Index(i).Interface()
			}
		} else {
			in <- i
		}
	}()
	return &defaultStream{in: in}
}
