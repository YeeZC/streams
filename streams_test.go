package streams

import (
	"fmt"
	"testing"
)

func TestFilter(t *testing.T) {
	Of([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).Filter(func(t interface{}) bool {
		return t.(int) > 5
	}).ForEach(func(i interface{}) {
		if i.(int) <= 5 {
			t.Fail()
		}
	})
}

func TestMap(t *testing.T) {
	Of([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).Map(func(t interface{}) interface{} {
		return t.(int) * 100
	}).ForEach(func(i interface{}) {
		fmt.Println(i)
	})
}

func TestDistinct(t *testing.T) {
	array := Of([]int{1, 2, 3, 3, 6, 6, 7, 8, 9, 10}).Distinct().ToArray().([]int)
	fmt.Printf("%v", array)
}

func TestSorted(t *testing.T) {
	opt := Of([]int{1, 2, 3, 3, 6, 6, 7, 8, 9, 10}).Distinct().Filter(func(t interface{}) bool {
		return t.(int) > 5
	}).Sorted(func(t1, t2 interface{}) int {
		return t2.(int) - t1.(int)
	}).FindAny()
	opt.IfPresent(func(t interface{}) {
		fmt.Println(t)
	})
}
