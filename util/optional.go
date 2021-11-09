package util

import "fmt"

var empty = Optional{}

type Optional struct {
	value interface{}
}

func (o Optional) Get() (interface{}, error) {
	if o.value != nil {
		return o.value, nil
	}
	return nil, fmt.Errorf("no value present")
}

func (o Optional) IsPresent() bool {
	return o.value != nil
}

func (o Optional) IfPresent(c Consumer) {
	if o.value != nil {
		c(o.value)
	}
}

func (o Optional) Filter(predicate Predicate) Optional {
	if !o.IsPresent() {
		return o
	}
	if predicate(o.value) {
		return o
	}
	return empty
}

func (o Optional) Map(function Function) Optional {
	if !o.IsPresent() {
		return o
	}
	return OfNullable(function(o.value))
}

func (o Optional) OrElse(v interface{}) interface{} {
	if o.value != nil {
		return o.value
	}
	return v
}

func Empty() Optional {
	return empty
}

func OfNullable(value interface{}) Optional {
	return Optional{value: value}
}
