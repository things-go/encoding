package form

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

//* encoder: custom type number/string slice/array

func EncodeCustomUnsignedIntegerSlice[T constraints.Unsigned](x any) ([]string, error) {
	return encodeCustomNumberSlice[T](x, func(u uint64) string {
		return strconv.FormatUint(u, 10)
	})
}

func EncodeCustomSignedIntegerSlice[T constraints.Unsigned](x any) ([]string, error) {
	return encodeCustomNumberSlice[T](x, func(u int64) string {
		return strconv.FormatInt(u, 10)
	})
}

func EncodeCustomFloat32Slice[T ~float32](x any) ([]string, error) {
	return encodeCustomNumberSlice[T](x, func(u float64) string {
		return strconv.FormatFloat(u, 'f', -1, 32)
	})
}

func EncodeCustomFloat64Slice[T ~float64](x any) ([]string, error) {
	return encodeCustomNumberSlice[T](x, func(u float64) string {
		return strconv.FormatFloat(u, 'f', -1, 64)
	})
}

func EncodeCustomStringSlice[T ~string](x any) ([]string, error) {
	vs, ok := x.([]T)
	if !ok {
		return nil, fmt.Errorf("")
	}
	has := false
	b := strings.Builder{}
	for _, vv := range vs {
		if vv == "" {
			continue
		}
		if has {
			b.WriteString(",")
		}
		has = true
		b.WriteString(string(vv))
	}
	return []string{b.String()}, nil
}

func encodeCustomNumberSlice[T constraints.Integer | constraints.Float, V uint64 | int64 | float64](x any, format func(V) string) ([]string, error) {
	vs, ok := x.([]T)
	if !ok {
		return nil, fmt.Errorf("")
	}
	has := false
	b := strings.Builder{}
	for _, v := range vs {
		vv := format(V(v))
		if vv == "" {
			continue
		}
		if has {
			b.WriteString(",")
		}
		has = true
		b.WriteString(vv)
	}
	return []string{b.String()}, nil
}
