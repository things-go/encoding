package form

import (
	"net/url"
	"testing"

	"github.com/go-playground/form/v4"
	"github.com/stretchr/testify/require"
)

type CustomInt int
type CustomInt8 int8
type CustomInt16 int16
type CustomInt32 int32
type CustomInt64 int64
type CustomUint uint
type CustomUint8 uint8
type CustomUint16 uint16
type CustomUint32 uint32
type CustomUint64 uint64
type CustomFloat32 float32
type CustomFloat64 float64
type CustomString string

type customDecode struct {
	I    []int           `json:"i"`
	I8   []int8          `json:"i8"`
	I16  []int16         `json:"i16"`
	I32  []int32         `json:"i32"`
	I64  []int64         `json:"i64"`
	U    []int           `json:"u"`
	U8   []int8          `json:"u8"`
	U16  []int16         `json:"u16"`
	U32  []int32         `json:"u32"`
	U64  []int64         `json:"u64"`
	F32  []float32       `json:"f32"`
	F64  []float64       `json:"f64"`
	S    []string        `json:"s"`
	Ci   []CustomInt     `json:"ci"`
	Ci8  []CustomInt8    `json:"ci8"`
	Ci16 []CustomInt16   `json:"ci16"`
	Ci32 []CustomInt32   `json:"ci32"`
	Ci64 []CustomInt64   `json:"ci64"`
	Cu   []CustomUint    `json:"cu"`
	Cu8  []CustomUint8   `json:"cu8"`
	Cu16 []CustomUint16  `json:"cu16"`
	Cu32 []CustomUint32  `json:"cu32"`
	Cu64 []CustomUint64  `json:"cu64"`
	Cf32 []CustomFloat32 `json:"cf32"`
	Cf64 []CustomFloat64 `json:"cf64"`
	Cs   []CustomString  `json:"cs"`
}

func Test_CustomTypeDecoder(t *testing.T) {
	dec := form.NewDecoder()
	dec.SetTagName("json")

	RegisterBuiltinSliceTypeDecoderComma(dec)
	dec.RegisterCustomTypeFunc(DecodeCustomIntSlice[CustomInt], []CustomInt{})
	dec.RegisterCustomTypeFunc(DecodeCustomInt8Slice[CustomInt8], []CustomInt8{})
	dec.RegisterCustomTypeFunc(DecodeCustomInt16Slice[CustomInt16], []CustomInt16{})
	dec.RegisterCustomTypeFunc(DecodeCustomInt32Slice[CustomInt32], []CustomInt32{})
	dec.RegisterCustomTypeFunc(DecodeCustomInt64Slice[CustomInt64], []CustomInt64{})
	dec.RegisterCustomTypeFunc(DecodeCustomUintSlice[CustomUint], []CustomUint{})
	dec.RegisterCustomTypeFunc(DecodeCustomUint8Slice[CustomUint8], []CustomUint8{})
	dec.RegisterCustomTypeFunc(DecodeCustomUint16Slice[CustomUint16], []CustomUint16{})
	dec.RegisterCustomTypeFunc(DecodeCustomUint32Slice[CustomUint32], []CustomUint32{})
	dec.RegisterCustomTypeFunc(DecodeCustomUint64Slice[CustomUint64], []CustomUint64{})
	dec.RegisterCustomTypeFunc(DecodeCustomFloat32Slice[CustomFloat32], []CustomFloat32{})
	dec.RegisterCustomTypeFunc(DecodeCustomFloat64Slice[CustomFloat64], []CustomFloat64{})
	dec.RegisterCustomTypeFunc(DecodeCustomStringSlice[CustomString], []CustomString{})
	got := customDecode{}
	err := dec.Decode(&got, url.Values{
		"i":    []string{"1,2,3", "4,5,6"},
		"i8":   []string{"81,82,83"},
		"i16":  []string{"161,162,163"},
		"i32":  []string{"321,322,323"},
		"i64":  []string{"641,642,643"},
		"u":    []string{"111111,222222,333333"},
		"u8":   []string{"11,22,33"},
		"u16":  []string{"1611,1622,1633"},
		"u32":  []string{"3211,3222,3233"},
		"u64":  []string{"6411,6422,6433"},
		"f32":  []string{"1.1,1.2,1.3"},
		"f64":  []string{"2.1,2.2,2.3"},
		"s":    []string{"a,b,c", "d,e,f"},
		"ci":   []string{"1,2,3"},
		"ci8":  []string{"81,82,83"},
		"ci16": []string{"161,162,163"},
		"ci32": []string{"321,322,323"},
		"ci64": []string{"641,642,643"},
		"cu":   []string{"111111,222222,333333"},
		"cu8":  []string{"11,22,33"},
		"cu16": []string{"1611,1622,1633"},
		"cu32": []string{"3211,3222,3233"},
		"cu64": []string{"6411,6422,6433"},
		"cf32": []string{"1.1,1.2,1.3"},
		"cf64": []string{"2.1,2.2,2.3"},
		"cs":   []string{"a,b,c"},
	})
	require.NoError(t, err)
	t.Logf("%#v", got)
}
