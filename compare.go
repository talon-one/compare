package compare

import (
	"reflect"
	"unsafe"

	forceexport "github.com/alangpierce/go-forceexport"
)

type Equaler interface {
	Equal(interface{}) bool
}

func Equal(x, y interface{}) bool {
	if x == nil || y == nil {
		return x == y
	}
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	if v1.Type() != v2.Type() {
		return false
	}
	return deepValueEqual(v1, v2, make(map[visit]bool), 0)
}

var reflectValueInterface func(v reflect.Value, safe bool) interface{}

func init() {
	err := forceexport.GetFunc(&reflectValueInterface, "reflect.valueInterface")
	if err != nil {
		panic(err)
	}
}

// During deepValueEqual, must keep track of checks that are
// in progress. The comparison algorithm assumes that all
// checks in progress are true when it reencounters them.
// Visited comparisons are stored in a map indexed by visit.
type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

// Tests for deep equality using reflected types. The map argument tracks
// comparisons that have already been seen, which allows short circuiting on
// recursive types.
func deepValueEqual(v1, v2 reflect.Value, visited map[visit]bool, depth int) bool {
	if !v1.IsValid() || !v2.IsValid() {
		return v1.IsValid() == v2.IsValid()
	}
	if v1.Type() != v2.Type() {
		return false
	}

	// if depth > 10 { panic("deepValueEqual") }	// for debugging

	// We want to avoid putting more in the visited map than we need to.
	// For any possible reference cycle that might be encountered,
	// hard(t) needs to return true for at least one of the types in the cycle.
	hard := func(k reflect.Kind) bool {
		switch k {
		case reflect.Map, reflect.Slice, reflect.Ptr, reflect.Interface:
			return true
		}
		return false
	}

	if v1.CanAddr() && v2.CanAddr() && hard(v1.Kind()) {
		addr1 := unsafe.Pointer(v1.UnsafeAddr())
		addr2 := unsafe.Pointer(v2.UnsafeAddr())
		if uintptr(addr1) > uintptr(addr2) {
			// Canonicalize order to reduce number of entries in visited.
			// Assumes non-moving garbage collector.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are already seen.
		typ := v1.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return true
		}

		// Remember for later.
		visited[v] = true
	}

	if i1 := reflectValueInterface(v1, false); i1 != nil {
		if i2 := reflectValueInterface(v2, false); i2 != nil {
			if eq, ok := i1.(Equaler); ok {
				return eq.Equal(i2)
			} else if eq, ok := i2.(Equaler); ok {
				return eq.Equal(i1)
			}
		}
	}

	switch v1.Kind() {
	case reflect.Array:
		for i := 0; i < v1.Len(); i++ {
			if !deepValueEqual(v1.Index(i), v2.Index(i), visited, depth+1) {
				return false
			}
		}
		return true
	case reflect.Slice:
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for i := 0; i < v1.Len(); i++ {
			if !deepValueEqual(v1.Index(i), v2.Index(i), visited, depth+1) {
				return false
			}
		}
		return true
	case reflect.Interface:
		if v1.IsNil() || v2.IsNil() {
			return v1.IsNil() == v2.IsNil()
		}
		return deepValueEqual(v1.Elem(), v2.Elem(), visited, depth+1)
	case reflect.Ptr:
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		return deepValueEqual(v1.Elem(), v2.Elem(), visited, depth+1)
	case reflect.Struct:
		for i, n := 0, v1.NumField(); i < n; i++ {
			if !deepValueEqual(v1.Field(i), v2.Field(i), visited, depth+1) {
				return false
			}
		}
		return true
	case reflect.Map:
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for _, k := range v1.MapKeys() {
			val1 := v1.MapIndex(k)
			val2 := v2.MapIndex(k)
			if !val1.IsValid() || !val2.IsValid() || !deepValueEqual(v1.MapIndex(k), v2.MapIndex(k), visited, depth+1) {
				return false
			}
		}
		return true
	case reflect.Func:
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		// Can't do better than this:
		return false
	default:
		// Normal equality suffices
		return reflectValueInterface(v1, false) == reflectValueInterface(v2, false)
	}
}
