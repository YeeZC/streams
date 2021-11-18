package filters

import "github.com/yeezc/streams/util"

var NotNil util.Predicate = func(i interface{}) bool {
	return i != nil
}
