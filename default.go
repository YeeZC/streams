package streams

import (
	"reflect"
	"sort"

	"github.com/yeezc/streams/util"
	"github.com/yeezc/streams/util/slices"
)

type defaultStream struct {
	in chan interface{}
}

func (s *defaultStream) Filter(predicate Predicate) Stream {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for elem := range s.in {
			if predicate(elem) {
				out <- elem
			}
		}
	}()
	return &defaultStream{in: out}
}

func (s *defaultStream) Map(function Function) Stream {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for elem := range s.in {
			out <- function(elem)
		}
	}()
	return &defaultStream{in: out}
}

func (s *defaultStream) FindAny() util.Optional {
	if elem, ok := <-s.in; ok {
		return util.OfNullable(elem)
	}
	return util.Empty()
}

func (s *defaultStream) Distinct() Stream {
	out := make(chan interface{})
	go func() {
		defer close(out)
		elems := make([]interface{}, 0)
		for elem := range s.in {
			if slices.Contains(elems, elem) {
				continue
			}
			elems = append(elems, elem)
			out <- elem
		}
	}()
	return &defaultStream{in: out}
}

func (s *defaultStream) Sorted(c Comparator) Stream {
	out := make(chan interface{})
	go func() {
		defer close(out)
		elems := make([]interface{}, 0)
		for elem := range s.in {
			elems = append(elems, elem)
		}
		comparable := &comparable{elements: elems, comparator: c}
		sort.Sort(comparable)
		for _, elem := range comparable.elements {
			out <- elem
		}
	}()
	return &defaultStream{in: out}
}

func (s *defaultStream) ForEach(consumer Consumer) {
	for elem := range s.in {
		consumer(elem)
	}
}

func (s *defaultStream) Reduce(identity T, op BinaryOperator) R {
	for elem := range s.in {
		identity = op(identity, elem)
	}
	return identity
}

func (s *defaultStream) ToArray() interface{} {
	var (
		ret reflect.Value
		set bool
	)
	nilCount := 0
	for elem := range s.in {
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
	return ret.Interface()
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
