package streams

import (
	"reflect"
	"sync"

	"github.com/reugn/go-streams"
	ext "github.com/reugn/go-streams/extension"
	"github.com/yeezc/streams/collectors"
	"github.com/yeezc/streams/util"
	"github.com/yeezc/streams/util/maps"
)

type Via interface {
	Via(streams.Flow) streams.Flow
	Out() <-chan interface{}
}

type Stream interface {
	Filter(predicate util.Predicate) Stream
	Map(function util.Function) Stream
	MapBreakable(function util.BreakableFunction) Stream
	FindAny() util.Optional
	Distinct() Stream
	Sorted(c util.Comparator) Stream
	Parallel(cnt uint) Stream
	Reverse() Stream
	ForEach(consumer util.Consumer)
	ForEachBreakable(function util.BreakableConsumer)
	Reduce(identity interface{}, op util.BinaryOperator) interface{}
	ToArray() interface{}
	Collect(collectors.Collector) interface{}
	Sum() float64
	Avg() float64
}

func EmptyStream() Stream {
	c := make(chan interface{})
	close(c)
	return &stream{via: ext.NewChanSource(c), parallel: 1}
}

func Of(i interface{}) Stream {
	in := make(chan interface{}, 1)
	go func() {
		defer close(in)
		v := reflect.ValueOf(i)
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				in <- v.Index(i).Interface()
			}
		case reflect.Map:
			it := v.MapRange()
			for it.Next() {
				key := it.Key().Interface()
				value := it.Value().Interface()
				in <- maps.NewEntry(key, value)
			}
		default:
			if m, ok := i.(sync.Map); ok {
				m.Range(func(key, value interface{}) bool {
					in <- maps.NewEntry(key, value)
					return true
				})
			} else {
				in <- i
			}
		}
	}()
	return &stream{via: ext.NewChanSource(in), parallel: 1}
}
