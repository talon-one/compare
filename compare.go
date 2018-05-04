package compare

import "reflect"

type Equaler interface {
	Equal(interface{}) bool
}

type Comparer interface {
	Compare(interface{}) int
}

func Compare(x, y interface{}) int {
	if cmp, ok := x.(Comparer); ok {
		return cmp.Compare(y)
	}
	if reflect.DeepEqual(x, y) {
		return 0
	}
	return 1
}

func Equal(x, y interface{}) bool {
	if eq, ok := x.(Equaler); ok {
		return eq.Equal(y)
	}
	return reflect.DeepEqual(x, y)
}
