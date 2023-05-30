package csize

import (
	"reflect"
)

type visitor struct {
	seen map[uintptr]struct{}
}

func (v *visitor) sizeOfList(inp reflect.Value) (sum uint32, ok bool) {
	switch inp.Type().Elem().Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.Uintptr:
		sum = uint32(inp.Len() * int(inp.Type().Elem().Size()))
	default:
		for i := 0; i < inp.Len(); i++ {
			s, ok := v.sizeOf(inp)
			if !ok {
				return 0, false
			}
			sum += s
		}
	}

	sum += uint32((inp.Cap() - inp.Len()) * int(inp.Type().Elem().Size()))
	return sum, true
}

func (v *visitor) sizeOfStruct(inp reflect.Value) (sum uint32, ok bool) {
	for i, n := 0, inp.NumField(); i < n; i++ {
		s, ok := v.sizeOf(inp.Field(i))
		if !ok {
			return 0, false
		}
		sum += s
	}
	return sum, true
}

func (v *visitor) sizeOf(inp reflect.Value) (uint32, bool) {
	switch inp.Kind() {
	case reflect.Array:
		return v.sizeOfList(inp)
	case reflect.Slice:
		if _, seen := v.seen[inp.Pointer()]; seen {
			return 0, false
		}
		v.seen[inp.Pointer()] = struct{}{}
		return v.sizeOfList(inp)
	case reflect.Struct:
		return v.sizeOfStruct(inp)
	case reflect.String:
		return uint32(inp.Len() + 1), true
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64:
		return uint32(int(inp.Type().Size())), true
	default:
		return 0, false
	}

}

func Of(value any) (uint32, bool) {
	v := &visitor{seen: make(map[uintptr]struct{})}
	return v.sizeOf(reflect.Indirect(reflect.ValueOf(value)))
}

func OfMust(value any) (out uint32) {
	var ok bool
	if out, ok = Of(value); !ok {
		panic("csize: incompatible value")
	}
	return
}
