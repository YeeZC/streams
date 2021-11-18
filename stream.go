package streams

import (
	"math"
	"reflect"
	"sort"

	ext "github.com/reugn/go-streams/extension"
	"github.com/reugn/go-streams/flow"
	"github.com/spf13/cast"
	"github.com/yeezc/streams/collectors"
	"github.com/yeezc/streams/util"
	"github.com/yeezc/streams/util/slices"
)

type stream struct {
	via      Via
	parallel uint
}

func (s *stream) Filter(predicate util.Predicate) Stream {
	via := s.via.Via(flow.NewFilter(flow.FilterFunc(predicate), s.parallel))
	return &stream{via: via, parallel: s.parallel}
}

func (s *stream) Map(function util.Function) Stream {
	via := s.via.Via(flow.NewMap(flow.MapFunc(function), s.parallel))
	return &stream{via: via, parallel: s.parallel}
}

func (s *stream) MapBreakable(function util.BreakableFunction) Stream {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for elem := range s.via.Out() {
			if inf, ok := function(elem); ok {
				out <- inf
			} else {
				break
			}
		}
	}()
	return &stream{via: ext.NewChanSource(out), parallel: s.parallel}
}

func (s *stream) FindAny() util.Optional {
	if elem, ok := <-s.via.Out(); ok {
		return util.OfNullable(elem)
	}
	return util.Empty()
}

func (s *stream) Distinct() Stream {
	out := make(chan interface{})
	go func() {
		defer close(out)
		elems := make([]interface{}, 0)
		for elem := range s.via.Out() {
			if !slices.Contains(elems, elem) {
				elems = append(elems, elem)
				out <- elem
			}
		}
	}()
	return &stream{via: ext.NewChanSource(out), parallel: s.parallel}
}

func (s *stream) Sorted(c util.Comparator) Stream {
	out := make(chan interface{})
	go func() {
		defer close(out)
		elems := make([]interface{}, 0)
		for elem := range s.via.Out() {
			elems = append(elems, elem)
		}
		comparable := &comparable{elements: elems, comparator: c}
		sort.Sort(comparable)
		for _, elem := range comparable.elements {
			out <- elem
		}
	}()
	return &stream{via: ext.NewChanSource(out), parallel: s.parallel}
}

func (s *stream) Parallel(cnt uint) Stream {
	s.parallel = cnt
	return s
}

func (s *stream) Reverse() Stream {
	out := make(chan interface{})
	go func() {
		defer close(out)
		elems := make([]interface{}, 0)
		for elem := range s.via.Out() {
			elems = append(elems, elem)
		}
		for i := len(elems) - 1; i >= 0; i-- {
			out <- elems[i]
		}
	}()
	return &stream{via: ext.NewChanSource(out), parallel: s.parallel}
}

func (s *stream) ForEach(consumer util.Consumer) {
	for elem := range s.via.Out() {
		consumer(elem)
	}
}

func (s *stream) ForEachBreakable(consumer util.BreakableConsumer) {
	for elem := range s.via.Out() {
		if ok := consumer(elem); !ok {
			break
		}
	}
}

func (s *stream) Reduce(identity interface{}, op util.BinaryOperator) interface{} {
	for elem := range s.via.Out() {
		identity = op(identity, elem)
	}
	return identity
}

func (s *stream) ToArray() interface{} {
	return s.Collect(collectors.Of(func() interface{} {
		return nil
	}, func(i1, i2 interface{}) interface{} {
		item := reflect.ValueOf(i2)
		if i1 == nil && i2 == nil {
			res := make([]interface{}, 0, 1)
			res = append(res, nil)
			return reflect.TypeOf(res)
		} else if i1 == nil {
			res := reflect.MakeSlice(reflect.SliceOf(item.Type()), 0, 10)
			return reflect.Append(res, item)
		}
		return reflect.Append(i1.(reflect.Value), item)
	}, func(i interface{}) interface{} {
		if v, ok := i.(reflect.Value); ok {
			return v.Interface()
		}
		return nil
	}))
}

func (s *stream) Collect(c collectors.Collector) interface{} {
	res := c.Supplier()()
	combiner := c.Combiner()
	for elem := range s.via.Out() {
		res = combiner(res, elem)
	}
	return c.Finisher()(res)
}

func (s *stream) Sum() float64 {
	return cast.ToFloat64(s.Reduce(float64(0), func(i1, i2 interface{}) interface{} {
		return i1.(float64) + cast.ToFloat64(i2)
	}))
}

func (s *stream) Avg() float64 {
	return s.Collect(&avg{}).(float64)
}

type comparable struct {
	comparator util.Comparator
	elements   []interface{}
}

func (c *comparable) Len() int {
	return len(c.elements)
}

func (c *comparable) Less(i, j int) bool {
	return c.comparator(c.elements[i], c.elements[j]) < 0
}

func (c *comparable) Swap(i, j int) {
	tmp := c.elements[i]
	c.elements[i] = c.elements[j]
	c.elements[j] = tmp
}

type avg struct {
	sum   float64
	count float64
}

func (c *avg) Supplier() util.Supplier {
	return func() interface{} {
		return nil
	}
}
func (c *avg) Combiner() util.BinaryOperator {
	return func(i1, i2 interface{}) interface{} {
		c.count++
		c.sum += cast.ToFloat64(i2)
		return nil
	}
}
func (c *avg) Finisher() util.Function {
	return func(i interface{}) interface{} {
		if c.count == 0 {
			return math.NaN()
		}
		return c.sum / c.count
	}
}
