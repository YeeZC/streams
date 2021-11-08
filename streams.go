package streams

import (
	"reflect"

	"github.com/reugn/go-streams"
	ext "github.com/reugn/go-streams/extension"
	"github.com/yeezc/streams/util"
)

type Via interface {
	Via(streams.Flow) streams.Flow
	Out() <-chan interface{}
}

type Stream interface {
	Filter(predicate util.Predicate) Stream
	Map(function util.Function) Stream
	FindAny() util.Optional
	Distinct() Stream
	Sorted(c util.Comparator) Stream
	Parallel(cnt uint) Stream
	ForEach(consumer util.Consumer)
	Reduce(identity interface{}, op util.BinaryOperator) interface{}
	ToArray() interface{}
}

func Empty() Stream {
	c := make(chan interface{})
	close(c)
	return &defaultStream{via: ext.NewChanSource(c), parallel: 1}
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
	return &defaultStream{via: ext.NewChanSource(in), parallel: 1}
}
