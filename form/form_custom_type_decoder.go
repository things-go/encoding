package form

import (
	"strconv"
	"strings"

	"github.com/go-playground/form/v4"
	"golang.org/x/exp/constraints"
)

func RegisterBuiltinSliceTypeDecoderComma(dec *form.Decoder) {
	dec.RegisterCustomTypeFunc(DecodeCustomIntSlice[int], []int{})
	dec.RegisterCustomTypeFunc(DecodeCustomInt8Slice[int8], []int8{})
	dec.RegisterCustomTypeFunc(DecodeCustomInt16Slice[int16], []int16{})
	dec.RegisterCustomTypeFunc(DecodeCustomInt32Slice[int32], []int32{})
	dec.RegisterCustomTypeFunc(DecodeCustomInt64Slice[int64], []int64{})
	dec.RegisterCustomTypeFunc(DecodeCustomUintSlice[uint], []uint{})
	dec.RegisterCustomTypeFunc(DecodeCustomUint8Slice[uint8], []uint8{})
	dec.RegisterCustomTypeFunc(DecodeCustomUint16Slice[uint16], []uint16{})
	dec.RegisterCustomTypeFunc(DecodeCustomUint32Slice[uint32], []uint32{})
	dec.RegisterCustomTypeFunc(DecodeCustomUint64Slice[uint64], []uint64{})
	dec.RegisterCustomTypeFunc(DecodeCustomFloat32Slice[float32], []float32{})
	dec.RegisterCustomTypeFunc(DecodeCustomFloat64Slice[float64], []float64{})
	dec.RegisterCustomTypeFunc(DecodeCustomStringSlice[string], []string{})
}

//* encoder: custom type number/string slice/array

func DecodeCustomUintSlice[T ~uint](values []string) (any, error) {
	return decodeNumber[T](values, func(s string) (uint64, error) {
		return strconv.ParseUint(s, 10, 0)
	})
}

func DecodeCustomUint8Slice[T ~uint8](values []string) (any, error) {
	return decodeNumber[T](values, func(s string) (uint64, error) {
		return strconv.ParseUint(s, 10, 8)
	})
}

func DecodeCustomUint16Slice[T ~uint16](values []string) (any, error) {
	return decodeNumber[T](values, func(s string) (uint64, error) {
		return strconv.ParseUint(s, 10, 16)
	})
}

func DecodeCustomUint32Slice[T ~uint32](values []string) (any, error) {
	return decodeNumber[T](values, func(s string) (uint64, error) {
		return strconv.ParseUint(s, 10, 32)
	})
}

func DecodeCustomUint64Slice[T ~uint64](values []string) (any, error) {
	return decodeNumber[T](values, func(s string) (uint64, error) {
		return strconv.ParseUint(s, 10, 64)
	})
}

func DecodeCustomIntSlice[T ~int](values []string) (any, error) {
	return decodeNumber[T](values, func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 0)
	})
}

func DecodeCustomInt8Slice[T ~int8](values []string) (any, error) {
	return decodeNumber[T](values, func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 8)
	})
}

func DecodeCustomInt16Slice[T ~int16](values []string) (any, error) {
	return decodeNumber[T](values, func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 16)
	})
}

func DecodeCustomInt32Slice[T ~int32](values []string) (any, error) {
	return decodeNumber[T](values, func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 32)
	})
}

func DecodeCustomInt64Slice[T ~int64](values []string) (any, error) {
	return decodeNumber[T](values, func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 64)
	})
}

func DecodeCustomFloat64Slice[T ~float64](values []string) (any, error) {
	return decodeNumber[T](values, func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	})
}

func DecodeCustomFloat32Slice[T ~float32](values []string) (any, error) {
	return decodeNumber[T](values, func(s string) (float64, error) {
		return strconv.ParseFloat(s, 32)
	})
}

func DecodeCustomStringSlice[T ~string](values []string) (any, error) {
	if len(values) == 0 {
		return []T{}, nil
	}
	// FIXME: make slice space
	ret := make([]T, 0)
	for _, s := range values {
		for _, v := range strings.Split(s, ",") {
			ret = append(ret, T(v))
		}
	}
	return ret, nil
}

func decodeNumber[T constraints.Integer | constraints.Float, V ~uint64 | ~int64 | ~float64](values []string, parse func(string) (V, error)) (any, error) {
	if len(values) == 0 {
		return []T{}, nil
	}
	// FIXME: make slice space
	ret := make([]T, 0)
	for _, s := range values {
		for _, v := range strings.Split(s, ",") {
			i, err := parse(v)
			if err != nil {
				return nil, err
			}
			ret = append(ret, T(i))
		}
	}
	return ret, nil
}
