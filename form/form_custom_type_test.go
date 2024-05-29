package form

import (
	"net/url"
	"testing"

	"github.com/go-playground/form/v4"
	"github.com/stretchr/testify/require"
)

type CustomBool bool
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
type CustomUintptr uintptr
type CustomFloat32 float32
type CustomFloat64 float64
type CustomString string

type customCodecValue struct {
	B    []bool          `json:"b"`
	I    []int           `json:"i"`
	I8   []int8          `json:"i8"`
	I16  []int16         `json:"i16"`
	I32  []int32         `json:"i32"`
	I64  []int64         `json:"i64"`
	U    []uint          `json:"u"`
	U8   []uint8         `json:"u8"`
	U16  []uint16        `json:"u16"`
	U32  []uint32        `json:"u32"`
	U64  []uint64        `json:"u64"`
	Up   []uintptr       `json:"up"`
	F32  []float32       `json:"f32"`
	F64  []float64       `json:"f64"`
	S    []string        `json:"s"`
	Cb   []CustomBool    `json:"cb"`
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
	Cup  []CustomUintptr `json:"cup"`
	Cf32 []CustomFloat32 `json:"cf32"`
	Cf64 []CustomFloat64 `json:"cf64"`
	Cs   []CustomString  `json:"cs"`

	Bb     []byte   `json:"bb"`
	EmptyN []int    `json:"emptyN"`
	EmptyS []string `json:"emptyS"`
}

var testCodecValue = customCodecValue{
	B:      []bool{false, true, true, false},
	I:      []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
	I8:     []int8{81, 82, 83},
	I16:    []int16{161, 162, 163},
	I32:    []int32{321, 322, 323},
	I64:    []int64{641, 642, 643},
	U:      []uint{111111, 222222, 333333},
	U8:     []uint8{11, 22, 33},
	U16:    []uint16{1611, 1622, 1633},
	U32:    []uint32{3211, 3222, 3233},
	U64:    []uint64{6411, 6422, 6433},
	Up:     []uintptr{99991, 99992, 99993},
	F32:    []float32{1.1, 1.2, 1.3},
	F64:    []float64{2.1, 2.2, 2.3},
	S:      []string{"a", "b", "c", "d", "e", "f"},
	Cb:     []CustomBool{false, true, true, false},
	Ci:     []CustomInt{1, 2, 3},
	Ci8:    []CustomInt8{81, 82, 83},
	Ci16:   []CustomInt16{161, 162, 163},
	Ci32:   []CustomInt32{321, 322, 323},
	Ci64:   []CustomInt64{641, 642, 643},
	Cu:     []CustomUint{111111, 222222, 333333},
	Cu8:    []CustomUint8{11, 22, 33},
	Cu16:   []CustomUint16{1611, 1622, 1633},
	Cu32:   []CustomUint32{3211, 3222, 3233},
	Cu64:   []CustomUint64{6411, 6422, 6433},
	Cup:    []CustomUintptr{99991, 99992, 99993},
	Cf32:   []CustomFloat32{1.1, 1.2, 1.3},
	Cf64:   []CustomFloat64{2.1, 2.2, 2.3},
	Cs:     []CustomString{"a", "b", "c"},
	Bb:     []byte{66, 77, 88},
	EmptyN: []int{},
	EmptyS: []string{},
}

var testDecoderUrlValues = url.Values{
	"b":    []string{"0,1,true,false"},
	"i":    []string{"1,2,3", "4,5,6,7,8,9"},
	"i8":   []string{"81,82,83"},
	"i16":  []string{"161,162,163"},
	"i32":  []string{"321,322,323"},
	"i64":  []string{"641,642,643"},
	"u":    []string{"111111,222222,333333"},
	"u8":   []string{"11,22,33"},
	"u16":  []string{"1611,1622,1633"},
	"u32":  []string{"3211,3222,3233"},
	"u64":  []string{"6411,6422,6433"},
	"up":   []string{"99991", "99992", "99993"},
	"f32":  []string{"1.1,1.2,1.3"},
	"f64":  []string{"2.1,2.2,2.3"},
	"s":    []string{"a,b,c", "d,e,f"},
	"cb":   []string{"0,1,true,false"},
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
	"cup":  []string{"99991", "99992", "99993"},
	"cf32": []string{"1.1,1.2,1.3"},
	"cf64": []string{"2.1,2.2,2.3"},
	"cs":   []string{"a,b,c"},

	"bb":     []string{"66,77,88"},
	"emptyN": []string{},
	"emptyS": []string{},
}

var testEncoderUrlValues = url.Values{
	"b":    []string{"false,true,true,false"},
	"i":    []string{"1,2,3,4,5,6,7,8,9"},
	"i8":   []string{"81,82,83"},
	"i16":  []string{"161,162,163"},
	"i32":  []string{"321,322,323"},
	"i64":  []string{"641,642,643"},
	"u":    []string{"111111,222222,333333"},
	"u8":   []string{"11,22,33"},
	"u16":  []string{"1611,1622,1633"},
	"u32":  []string{"3211,3222,3233"},
	"u64":  []string{"6411,6422,6433"},
	"up":   []string{"99991,99992,99993"},
	"f32":  []string{"1.1,1.2,1.3"},
	"f64":  []string{"2.1,2.2,2.3"},
	"s":    []string{"a,b,c,d,e,f"},
	"cb":   []string{"false,true,true,false"},
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
	"cup":  []string{"99991,99992,99993"},
	"cf32": []string{"1.1,1.2,1.3"},
	"cf64": []string{"2.1,2.2,2.3"},
	"cs":   []string{"a,b,c"},

	"bb":     []string{"66,77,88"},
	"emptyN": []string{},
	"emptyS": []string{},
}

func Test_CustomType_Encoder(t *testing.T) {
	enc := form.NewEncoder()
	enc.SetTagName("json")

	RegisterBuiltinTypeEncoderSliceToCommaString(enc)
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[CustomBool](), []CustomBool{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[CustomInt](), []CustomInt{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[CustomInt8](), []CustomInt8{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[CustomInt16](), []CustomInt16{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[CustomInt32](), []CustomInt32{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[CustomInt64](), []CustomInt64{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[CustomUint](), []CustomUint{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[CustomUint8](), []CustomUint8{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[CustomUint16](), []CustomUint16{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[CustomUint32](), []CustomUint32{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[CustomUint64](), []CustomUint64{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[CustomUintptr](), []CustomUintptr{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[CustomFloat32](), []CustomFloat32{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[CustomFloat64](), []CustomFloat64{})
	enc.RegisterCustomTypeFunc(EncodeSlice2CommaString[CustomString](), []CustomString{})

	t.Run("valid type", func(t *testing.T) {
		got, err := enc.Encode(testCodecValue)
		require.NoError(t, err)
		require.Equal(t, testEncoderUrlValues, got)
	})

	t.Run("invalid type which no register", func(t *testing.T) {
		type CustomXXX int
		type Custom struct {
			A []CustomXXX `json:"a"`
		}

		got, err := enc.Encode(Custom{
			A: []CustomXXX{1, 2},
		})
		require.NoError(t, err)
		// want behaviour, but not
		require.NotEqual(t, url.Values{"a": []string{"1,2"}}, got)
		// default behaviour
		require.Equal(t, url.Values{"a": []string{"1", "2"}}, got)
	})
}

func Test_CustomType_Decoder(t *testing.T) {
	dec := form.NewDecoder()
	dec.SetTagName("json")

	RegisterBuiltinTypeDecoderCommaStringToSlice(dec)
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[CustomBool](), []CustomBool{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[CustomInt](), []CustomInt{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[CustomInt8](), []CustomInt8{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[CustomInt16](), []CustomInt16{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[CustomInt32](), []CustomInt32{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[CustomInt64](), []CustomInt64{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[CustomUint](), []CustomUint{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[CustomUint8](), []CustomUint8{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[CustomUint16](), []CustomUint16{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[CustomUint32](), []CustomUint32{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[CustomUint64](), []CustomUint64{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[CustomUintptr](), []CustomUintptr{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[CustomFloat32](), []CustomFloat32{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[CustomFloat64](), []CustomFloat64{})
	dec.RegisterCustomTypeFunc(DecodeCommaString2Slice[CustomString](), []CustomString{})

	t.Run("valid type", func(t *testing.T) {
		got := customCodecValue{}

		err := dec.Decode(&got, testDecoderUrlValues)
		require.NoError(t, err)
		require.Equal(t, testCodecValue, got)
	})

	t.Run("invalid type which no register", func(t *testing.T) {
		type CustomXXX int
		type Custom struct {
			A []CustomXXX `json:"a"`
		}
		got := Custom{}
		err := dec.Decode(&got, url.Values{
			"a": []string{"1,2"},
		})
		require.Error(t, err)
		t.Log(err)
	})
	t.Run("skip if empty string and number type", func(t *testing.T) {
		type Custom struct {
			A []int `json:"a"`
		}
		got1 := Custom{}
		err := dec.Decode(&got1, url.Values{
			"a": []string{""},
		})
		require.NoError(t, err)
		require.Equal(t, Custom{A: []int{}}, got1)

		got2 := Custom{}
		err = dec.Decode(&got2, url.Values{
			"a": []string{"1,,3"},
		})
		require.NoError(t, err)
		require.Equal(t, Custom{A: []int{1, 3}}, got2)
	})

	t.Run("invalid value", func(t *testing.T) {
		got := customCodecValue{}
		err := dec.Decode(&got, url.Values{
			"i": []string{"a,b"},
		})
		require.Error(t, err)
	})

}
