package streams

import (
	"fmt"
	"testing"

	"github.com/yeezc/streams/util"
)

func TestFilter(t *testing.T) {
	Of([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).Filter(func(t util.T) bool {
		return t.(int) > 5
	}).ForEach(func(i util.T) {
		if i.(int) <= 5 {
			t.Fail()
		}
	})
}

func TestMap(t *testing.T) {
	Of([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).Map(func(t util.T) util.R {
		return t.(int) * 100
	}).ForEach(func(i util.T) {
		fmt.Println(i)
	})
}

func TestDistinct(t *testing.T) {
	array := Of([]int{1, 2, 3, 3, 6, 6, 7, 8, 9, 10}).Distinct().ToArray().([]int)
	fmt.Printf("%v", array)
}

func TestSorted(t *testing.T) {
	opt := Of([]int{1, 2, 3, 3, 6, 6, 7, 8, 9, 10}).Distinct().Filter(func(t util.T) bool {
		return t.(int) > 5
	}).Sorted(func(t1, t2 util.T) int {
		return t2.(int) - t1.(int)
	}).FindAny()
	opt.IfPresent(func(t util.T) {
		fmt.Println(t)
	})
}
