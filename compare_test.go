package compare

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type IntStruct struct {
	Integer int
}

func (i IntStruct) Equal(v interface{}) bool {
	switch x := v.(type) {
	case IntStruct:
		return i.Integer == x.Integer
	case *IntStruct:
		return i.Integer == x.Integer
	}
	return false
}

func (i IntStruct) Compare(v interface{}) int {
	switch x := v.(type) {
	case IntStruct:
		if i.Integer == x.Integer {
			return 0
		} else if i.Integer > x.Integer {
			return 1
		}
		return -1
	case *IntStruct:
		if i.Integer == x.Integer {
			return 0
		} else if i.Integer > x.Integer {
			return 1
		}
		return -1
	}
	return 1
}

type StringStruct struct {
	String string
}

func (i StringStruct) Equal(v interface{}) bool {
	switch x := v.(type) {
	case StringStruct:
		return i.String == x.String
	case *StringStruct:
		return i.String == x.String
	}
	return false
}

func TestEqual(t *testing.T) {
	s := struct{}{}
	x := IntStruct{1}
	y := IntStruct{2}
	z := StringStruct{"Hello"}

	require.Equal(t, true, Equal(x, x))
	require.Equal(t, true, Equal(x, &x))
	require.Equal(t, true, Equal(&x, &x))

	require.Equal(t, false, Equal(x, y))
	require.Equal(t, false, Equal(x, &y))
	require.Equal(t, false, Equal(&x, &y))

	require.Equal(t, false, Equal(x, z))
	require.Equal(t, false, Equal(x, &z))
	require.Equal(t, false, Equal(&x, &z))

	require.Equal(t, false, Equal(x, s))
	require.Equal(t, false, Equal(x, &s))
	require.Equal(t, false, Equal(&x, &s))

	require.Equal(t, false, Equal(s, x))
	require.Equal(t, false, Equal(s, &x))
	require.Equal(t, false, Equal(&s, &x))
}

func TestCompare(t *testing.T) {
	s := struct{}{}
	v := IntStruct{1}
	x := IntStruct{2}
	y := IntStruct{3}
	z := StringStruct{"Hello"}

	require.Equal(t, 0, Compare(v, v))
	require.Equal(t, 0, Compare(v, &v))
	require.Equal(t, 0, Compare(&v, &v))

	require.Equal(t, 1, Compare(x, v))
	require.Equal(t, 1, Compare(x, &v))
	require.Equal(t, 1, Compare(&x, &v))

	require.Equal(t, -1, Compare(x, y))
	require.Equal(t, -1, Compare(x, &y))
	require.Equal(t, -1, Compare(&x, &y))

	require.Equal(t, 1, Compare(v, s))
	require.Equal(t, 1, Compare(v, &s))
	require.Equal(t, 1, Compare(&v, &s))

	require.Equal(t, 1, Compare(s, v))
	require.Equal(t, 1, Compare(s, &v))
	require.Equal(t, 1, Compare(&s, &v))

	require.Equal(t, 1, Compare(v, z))
	require.Equal(t, 1, Compare(v, &z))
	require.Equal(t, 1, Compare(&v, &z))

	s1 := struct{ A int }{7}
	s2 := struct{ A int }{7}

	require.Equal(t, 0, Compare(s1, s2))
}
