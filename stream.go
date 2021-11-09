package streams

import (
	"reflect"
	"sort"

	ext "github.com/reugn/go-streams/extension"
	"github.com/reugn/go-streams/flow"
	. "github.com/yeezc/streams/util"
	"github.com/yeezc/streams/util/slices"
)

type stream struct {
	via      Via
	parallel uint
}

func (s *stream) Filter(predicate Predicate) Stream {
	via := s.via.Via(flow.NewFilter(flow.FilterFunc(predicate), s.parallel))
	return &stream{via: via, parallel: s.parallel}
}

func (s *stream) Map(function Function) Stream {
	via := s.via.Via(flow.NewMap(flow.MapFunc(function), s.parallel))
	return &stream{via: via, parallel: s.parallel}
}

func (s *stream) FindAny() Optional {
	if elem, ok := <-s.via.Out(); ok {
		return OfNullable(elem)
	}
	return Empty()
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

func (s *stream) Sorted(c Comparator) Stream {
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

func (s *stream) ForEach(consumer Consumer) {
	for elem := range s.via.Out() {
		consumer(elem)
	}
}

func (s *stream) Reduce(identity interface{}, op BinaryOperator) interface{} {
	for elem := range s.via.Out() {
		identity = op(identity, elem)
	}
	return identity
}

func (s *stream) ToArray() interface{} {
	var (
		ret reflect.Value
		set bool
	)
	nilCount := 0
	for elem := range s.via.Out() {
		el := reflect.ValueOf(elem)
		if !set && elem != nil {
			ret = reflect.MakeSlice(reflect.SliceOf(el.Type()), 0, nilCount)
			for i := 0; i < nilCount; i++ {
				reflect.Append(ret, reflect.ValueOf(nil))
			}
			ret = reflect.Append(ret, el)
			nilCount = 0
			set = true
		} else if !set && elem == nil {
			nilCount++
		} else {
			ret = reflect.Append(ret, el)
		}
	}
	if ret.CanInterface() {
		return ret.Interface()
	}
	return nil
}

type comparable struct {
	comparator Comparator
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
