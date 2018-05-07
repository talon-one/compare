package compare

import (
	"testing"
	"time"

	"math/rand"

	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type IntStruct struct {
	Integer int
	rnd     int32
}

func NewIntStruct(i int) IntStruct {
	return IntStruct{
		Integer: i,
		rnd:     rand.Int31(),
	}
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

type StringStruct struct {
	String string
	rnd    int32
}

func NewStringStruct(s string) StringStruct {
	return StringStruct{
		String: s,
		rnd:    rand.Int31(),
	}
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
	x := NewIntStruct(1)
	y := NewIntStruct(2)
	z := NewStringStruct("Hello")

	require.Equal(t, true, Equal(x, x))
	require.Equal(t, false, Equal(x, &x))
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

func TestEqualSlice(t *testing.T) {
	x := []interface{}{
		NewIntStruct(1),
		NewIntStruct(2),
	}
	y := []interface{}{
		NewIntStruct(2),
		NewIntStruct(2),
	}
	z := []interface{}{
		NewIntStruct(1),
		NewStringStruct("Hello"),
	}
	require.Equal(t, true, Equal(x, x))
	require.Equal(t, false, Equal(x, y))
	require.Equal(t, false, Equal(x, z))
}

func TestEqualMap(t *testing.T) {
	x := map[string]interface{}{
		"A": NewIntStruct(1),
		"B": NewIntStruct(2),
	}
	y := map[string]interface{}{
		"A": NewIntStruct(2),
		"B": NewIntStruct(2),
	}
	z := map[string]interface{}{
		"A": NewIntStruct(1),
		"B": NewStringStruct("Hello"),
	}
	require.Equal(t, true, Equal(x, x))
	require.Equal(t, false, Equal(x, y))
	require.Equal(t, false, Equal(x, z))
}

func TestEqualPrivate(t *testing.T) {
	type St1 struct {
		str string
	}
	require.Equal(t, true, Equal(St1{"Hello"}, St1{"Hello"}))

	type St2 struct {
		i IntStruct
	}
	require.Equal(t, true, Equal(St2{NewIntStruct(2)}, St2{NewIntStruct(2)}))
}
