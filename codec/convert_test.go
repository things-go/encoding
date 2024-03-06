package codec

import (
	"reflect"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func protoEqual(a, b any) bool {
	return proto.Equal(a.(proto.Message), b.(proto.Message))
}

type testStruct[T any] struct {
	name    string
	input   string
	output  T
	wantErr bool
}

func test_BuiltinType[T any](t *testing.T, tests []testStruct[T], checkFunc func(string) (T, error), equal func(any, any) bool) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, err := checkFunc(tt.input)
			switch {
			case err != nil && !tt.wantErr:
				t.Errorf("got unexpected error\n%#v", err)
			case err == nil && tt.wantErr:
				t.Errorf("did not error when expected")
			case !equal(ts, tt.output):
				t.Errorf(
					"when testing %s; got\n%#v\nexpected\n%#v",
					tt.name,
					ts,
					tt.output,
				)
			}
		})
	}
}
func test_BuiltinTypeSlices[T any](t *testing.T, tests []testStruct[T], checkFunc func(string, string) (T, error), equal func(any, any) bool) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, err := checkFunc(tt.input, ",")
			switch {
			case err != nil && !tt.wantErr:
				t.Errorf("got unexpected error\n%#v", err)
			case err == nil && tt.wantErr:
				t.Errorf("did not error when expected")
			case !equal(ts, tt.output):
				t.Errorf(
					"when testing %s; got\n%#v\nexpected\n%#v",
					tt.name,
					ts,
					tt.output,
				)
			}
		})
	}
}

func Test_String(t *testing.T) {
	tests := []testStruct[string]{
		{
			name:    "",
			input:   "hello world",
			output:  "hello world",
			wantErr: false,
		},
	}
	test_BuiltinType(t, tests, String, reflect.DeepEqual)
}

func Test_StringSlice(t *testing.T) {
	tests := []testStruct[[]string]{
		{
			name:    "",
			input:   "hello,world",
			output:  []string{"hello", "world"},
			wantErr: false,
		},
	}
	test_BuiltinTypeSlices(t, tests, StringSlice, reflect.DeepEqual)
}

func Test_Bool(t *testing.T) {
	tests := []testStruct[bool]{
		{
			name:    "true",
			input:   "true",
			output:  true,
			wantErr: false,
		},
		{
			name:    "T",
			input:   "true",
			output:  true,
			wantErr: false,
		},
		{
			name:    "invalid",
			input:   "11",
			output:  false,
			wantErr: true,
		},
	}
	test_BuiltinType(t, tests, Bool, reflect.DeepEqual)
}

func Test_BoolSlice(t *testing.T) {
	tests := []testStruct[[]bool]{
		{
			name:    "true|false",
			input:   "true,false",
			output:  []bool{true, false},
			wantErr: false,
		},
		{
			name:    "T|f",
			input:   "T,f",
			output:  []bool{true, false},
			wantErr: false,
		},
		{
			name:    "invalid",
			input:   "11",
			output:  nil,
			wantErr: true,
		},
	}
	test_BuiltinTypeSlices(t, tests, BoolSlice, reflect.DeepEqual)
}

func Test_Float64(t *testing.T) {
	tests := []testStruct[float64]{
		{
			name:    "true",
			input:   "1.0",
			output:  1.0,
			wantErr: false,
		},
		{
			name:    "invalid",
			input:   "1.x",
			output:  0,
			wantErr: true,
		},
	}
	test_BuiltinType(t, tests, Float64, reflect.DeepEqual)
}

func Test_Float64Slice(t *testing.T) {
	tests := []testStruct[[]float64]{
		{
			name:    "true",
			input:   "1.0,1.1",
			output:  []float64{1.0, 1.1},
			wantErr: false,
		},
		{
			name:    "invalid",
			input:   "1.x,x",
			output:  nil,
			wantErr: true,
		},
	}
	test_BuiltinTypeSlices(t, tests, Float64Slice, reflect.DeepEqual)
}

func Test_Float32(t *testing.T) {
	tests := []testStruct[float32]{
		{
			name:    "true",
			input:   "1.0",
			output:  1.0,
			wantErr: false,
		},
		{
			name:    "invalid",
			input:   "1.x",
			output:  0,
			wantErr: true,
		},
	}
	test_BuiltinType(t, tests, Float32, reflect.DeepEqual)
}

func Test_Float32Slice(t *testing.T) {
	tests := []testStruct[[]float32]{
		{
			name:    "true",
			input:   "1.0,1.1",
			output:  []float32{1.0, 1.1},
			wantErr: false,
		},
		{
			name:    "invalid",
			input:   "1.x,x",
			output:  nil,
			wantErr: true,
		},
	}
	test_BuiltinTypeSlices(t, tests, Float32Slice, reflect.DeepEqual)
}

func Test_Timestamp(t *testing.T) {
	tests := []testStruct[*timestamppb.Timestamp]{
		{
			name:  "a valid RFC3339 timestamp",
			input: `"2016-05-10T10:19:13.123Z"`,
			output: &timestamppb.Timestamp{
				Seconds: 1462875553,
				Nanos:   123000000,
			},
			wantErr: false,
		},
		{
			name:  "a valid RFC3339 timestamp without double quotation",
			input: "2016-05-10T10:19:13.123Z",
			output: &timestamppb.Timestamp{
				Seconds: 1462875553,
				Nanos:   123000000,
			},
			wantErr: false,
		},
		{
			name:    "invalid timestamp",
			input:   `"05-10-2016T10:19:13.123Z"`,
			output:  nil,
			wantErr: true,
		},
		{
			name:    "JSON number",
			input:   "123",
			output:  nil,
			wantErr: true,
		},
		{
			name:    "JSON bool",
			input:   "true",
			output:  nil,
			wantErr: true,
		},
	}
	test_BuiltinType(t, tests, Timestamp, protoEqual)
}

func Test_Duration(t *testing.T) {
	tests := []testStruct[*durationpb.Duration]{
		{
			name:  "a valid duration",
			input: `"123.456s"`,
			output: &durationpb.Duration{
				Seconds: 123,
				Nanos:   456000000,
			},
			wantErr: false,
		},
		{
			name:  "a valid duration without double quotation",
			input: "123.456s",
			output: &durationpb.Duration{
				Seconds: 123,
				Nanos:   456000000,
			},
			wantErr: false,
		},
		{
			name:    "invalid duration",
			input:   `"123years"`,
			output:  nil,
			wantErr: true,
		},
		{
			name:    "JSON number",
			input:   "123",
			output:  nil,
			wantErr: true,
		},
		{
			name:    "JSON bool",
			input:   "true",
			output:  nil,
			wantErr: true,
		},
	}
	test_BuiltinType(t, tests, Duration, protoEqual)
}
