package form

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/form/v4"
	"golang.org/x/exp/constraints"
)

// * encode/decode: custom type number/string slice/array

// RegisterBuiltinTypeEncoderSliceToCommaString register to form.Encoder.
// encode a slice to a comma-separated string.
// NOTE: slice element type only support
//
// `bool`
// `int`, `int8`, `int16`, `int32`, `int64`
// `uint`, `uint8`, `uint16`, `uint32`, `uint64`
// `float32`, `float64`
// `string`, `uintptr`
func RegisterBuiltinTypeEncoderSliceToCommaString(enc *form.Encoder) {
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[bool](), []bool{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[int](), []int{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[int8](), []int8{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[int16](), []int16{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[int32](), []int32{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[int64](), []int64{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[uint](), []uint{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[uint8](), []uint8{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[uint16](), []uint16{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[uint32](), []uint32{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[uint64](), []uint64{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[float32](), []float32{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[float64](), []float64{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[string](), []string{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[uintptr](), []uintptr{})
}

// RegisterBuiltinTypeDecoderCommaStringToSlice register to form.Decoder.
// decode a comma-separated string to slice.
// NOTE: slice element type only support
//
// `bool`
// `int`, `int8`, `int16`, `int32`, `int64`
// `uint`, `uint8`, `uint16`, `uint32`, `uint64`
// `float32`, `float64`
// `string`, `uintptr`
func RegisterBuiltinTypeDecoderCommaStringToSlice(dec *form.Decoder) {
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[bool](), []bool{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[int](), []int{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[int8](), []int8{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[int16](), []int16{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[int32](), []int32{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[int64](), []int64{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[uint](), []uint{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[uint8](), []uint8{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[uint16](), []uint16{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[uint32](), []uint32{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[uint64](), []uint64{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[float32](), []float32{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[float64](), []float64{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[string](), []string{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[uintptr](), []uintptr{})
}

// EncodeSlice2CommaString encode a slice to a comma-separated string.
func EncodeSlice2CommaString[T constraints.Integer | constraints.Float | ~string | ~bool]() func(x any) ([]string, error) {
	t := reflect.TypeOf([]T{})
	return func(x any) ([]string, error) {
		return EncodeSliceToCommaString(t, x)
	}
}

// EncodeSliceToCommaString encode a slice to a comma-separated string.
// NOTE: slice element only support `constraints.Integer | constraints.Float | ~string | ~bool`
func EncodeSliceToCommaString(t reflect.Type, x any) ([]string, error) {
	if t.Kind() != reflect.Slice {
		return nil, &form.InvalidEncodeError{Type: t}
	}
	vx := reflect.ValueOf(x)
	if vx.Kind() != reflect.Slice {
		return nil, &form.InvalidEncodeError{Type: t}
	}

	teKind := t.Elem().Kind()
	if reflect.TypeOf(x).Elem().Kind() != teKind {
		return nil, &form.InvalidEncodeError{Type: t}
	}
	if vx.Len() == 0 {
		return []string{}, nil
	}

	has := false
	b := strings.Builder{}
	for i := 0; i < vx.Len(); i++ {
		fv := vx.Index(i)
		val := ""
		switch teKind {
		case reflect.Bool:
			val = strconv.FormatBool(fv.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val = strconv.FormatInt(fv.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val = strconv.FormatUint(fv.Uint(), 10)
		case reflect.Float32:
			val = strconv.FormatFloat(fv.Float(), 'f', -1, 32)
		case reflect.Float64:
			val = strconv.FormatFloat(fv.Float(), 'f', -1, 64)
		case reflect.String:
			val = fv.String()
		default:
			return nil, &form.InvalidEncodeError{Type: t}
		}
		if has {
			b.WriteString(",")
		}
		has = true
		b.WriteString(val)
	}
	return []string{b.String()}, nil
}

// DecodeCommaString2Slice decode a comma-separated string to a slice.
func DecodeCommaString2Slice[T constraints.Integer | constraints.Float | ~string | ~bool]() func(values []string) (any, error) {
	t := reflect.TypeOf([]T{})
	return func(values []string) (any, error) {
		return DecodeCommaString22Slice(t, values)
	}
}

// DecodeCommaString22Slice decode a comma-separated string to a slice.
// NOTE: slice element only support `constraints.Integer | constraints.Float | ~string | ~bool`
func DecodeCommaString22Slice(t reflect.Type, values []string) (any, error) {
	if t.Kind() != reflect.Slice {
		return nil, &form.InvalidDecoderError{Type: t}
	}
	if len(values) == 0 {
		return reflect.MakeSlice(t, 0, 0).Interface(), nil
	}

	ret := reflect.MakeSlice(t, 0, 8)
	te := t.Elem()
	teKind := te.Kind()
	for _, s := range values {
		elements := strings.Split(s, ",")
		if oldLen, oldCap := ret.Len(), ret.Cap(); oldCap < oldLen+len(elements) {
			newCap := growCap(oldCap, oldCap+len(elements))
			nret := reflect.MakeSlice(t, oldLen, newCap)
			reflect.Copy(nret, ret)
			ret = nret
		}
		for _, ss := range elements {
			val := reflect.New(te).Elem()
			switch teKind {
			case reflect.Bool:
				i, err := strconv.ParseBool(ss)
				if err != nil {
					return nil, err
				}
				val.SetBool(i)
			case reflect.Int:
				i, err := strconv.ParseInt(ss, 10, 0)
				if err != nil {
					return nil, err
				}
				val.SetInt(i)
			case reflect.Int8:
				i, err := strconv.ParseInt(ss, 10, 8)
				if err != nil {
					return nil, err
				}
				val.SetInt(i)
			case reflect.Int16:
				i, err := strconv.ParseInt(ss, 10, 16)
				if err != nil {
					return nil, err
				}
				val.SetInt(i)
			case reflect.Int32:
				i, err := strconv.ParseInt(ss, 10, 32)
				if err != nil {
					return nil, err
				}
				val.SetInt(i)
			case reflect.Int64:
				i, err := strconv.ParseInt(ss, 10, 64)
				if err != nil {
					return nil, err
				}
				val.SetInt(i)
			case reflect.Uint:
				i, err := strconv.ParseUint(ss, 10, 0)
				if err != nil {
					return nil, err
				}
				val.SetUint(i)
			case reflect.Uint8:
				i, err := strconv.ParseUint(ss, 10, 8)
				if err != nil {
					return nil, err
				}
				val.SetUint(i)
			case reflect.Uint16:
				i, err := strconv.ParseUint(ss, 10, 16)
				if err != nil {
					return nil, err
				}
				val.SetUint(i)
			case reflect.Uint32:
				i, err := strconv.ParseUint(ss, 10, 32)
				if err != nil {
					return nil, err
				}
				val.SetUint(i)
			case reflect.Uint64:
				i, err := strconv.ParseUint(ss, 10, 64)
				if err != nil {
					return nil, err
				}
				val.SetUint(i)
			case reflect.Float32:
				i, err := strconv.ParseFloat(ss, 32)
				if err != nil {
					return nil, err
				}
				val.SetFloat(i)
			case reflect.Float64:
				i, err := strconv.ParseFloat(ss, 64)
				if err != nil {
					return nil, err
				}
				val.SetFloat(i)
			case reflect.String:
				val.SetString(ss)
			default:
				return nil, &form.InvalidDecoderError{Type: t}
			}
			ret = reflect.Append(ret, val)
		}
	}
	return ret.Interface(), nil
}

func growCap(oldCap, cap int) int {
	newCap := oldCap
	doubleCap := newCap + newCap
	if cap > doubleCap {
		newCap = cap
	} else {
		const threshold = 256
		if oldCap < threshold {
			newCap = doubleCap
		} else {
			for 0 < newCap && newCap < cap {
				newCap += (newCap + 3*threshold) / 4
			}
		}
		if newCap <= 0 {
			newCap = cap
		}
	}
	return newCap
}
