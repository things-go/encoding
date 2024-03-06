package jsonpb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/things-go/encoding/codec"
)

// Codec is a Marshaler which marshals/unmarshals into/from JSON
// with the "google.golang.org/protobuf/encoding/protojson" marshaler.
// It supports the full functionality of protobuf unlike JSONBuiltin.
//
// The NewDecoder method returns a DecoderWrapper, so the underlying
// *json.Decoder methods can be used.
type Codec struct {
	protojson.MarshalOptions
	protojson.UnmarshalOptions
}

// ContentType always Returns "application/json; charset=utf-8".
func (*Codec) ContentType(_ any) string {
	return "application/json; charset=utf-8"
}

func (c *Codec) Marshal(v any) ([]byte, error) {
	if _, ok := v.(proto.Message); !ok {
		return c.marshalNonProtoField(v)
	}

	var buf bytes.Buffer
	if err := c.marshalTo(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Codec) marshalTo(w io.Writer, v any) error {
	p, ok := v.(proto.Message)
	if !ok {
		buf, err := c.marshalNonProtoField(v)
		if err != nil {
			return err
		}
		_, err = w.Write(buf)
		return err
	}
	b, err := c.MarshalOptions.Marshal(p)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}

var (
	// protoMessageType is stored to prevent constant lookup of the same type at runtime.
	protoMessageType = reflect.TypeOf((*proto.Message)(nil)).Elem()
)

// marshalNonProto marshals a non-message field of a protobuf message.
// This function does not correctly marshal arbitrary data structures into JSON,
// it is only capable of marshaling non-message field values of protobuf,
// i.e. primitive types, enums; pointers to primitives or enums; maps from
// integer/string types to primitives/enums/pointers to messages.
func (c *Codec) marshalNonProtoField(v any) ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return []byte("null"), nil
		}
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Slice {
		if rv.IsNil() {
			if c.EmitUnpopulated {
				return []byte("[]"), nil
			}
			return []byte("null"), nil
		}

		if rv.Type().Elem().Implements(protoMessageType) {
			var buf bytes.Buffer
			err := buf.WriteByte('[')
			if err != nil {
				return nil, err
			}
			for i := 0; i < rv.Len(); i++ {
				if i != 0 {
					err = buf.WriteByte(',')
					if err != nil {
						return nil, err
					}
				}
				if err = c.marshalTo(&buf, rv.Index(i).Interface().(proto.Message)); err != nil {
					return nil, err
				}
			}
			err = buf.WriteByte(']')
			if err != nil {
				return nil, err
			}

			return buf.Bytes(), nil
		}

		if rv.Type().Elem().Implements(typeProtoEnum) {
			var buf bytes.Buffer
			err := buf.WriteByte('[')
			if err != nil {
				return nil, err
			}
			for i := 0; i < rv.Len(); i++ {
				if i != 0 {
					err = buf.WriteByte(',')
					if err != nil {
						return nil, err
					}
				}
				if c.UseEnumNumbers {
					_, err = buf.WriteString(strconv.FormatInt(rv.Index(i).Int(), 10))
				} else {
					_, err = buf.WriteString("\"" + rv.Index(i).Interface().(protoEnum).String() + "\"")
				}
				if err != nil {
					return nil, err
				}
			}
			err = buf.WriteByte(']')
			if err != nil {
				return nil, err
			}

			return buf.Bytes(), nil
		}
	}

	if rv.Kind() == reflect.Map {
		m := make(map[string]*json.RawMessage)
		for _, k := range rv.MapKeys() {
			buf, err := c.Marshal(rv.MapIndex(k).Interface())
			if err != nil {
				return nil, err
			}
			m[fmt.Sprintf("%v", k.Interface())] = (*json.RawMessage)(&buf)
		}
		if c.Indent != "" {
			return json.MarshalIndent(m, "", c.Indent)
		}
		return json.Marshal(m)
	}
	if enum, ok := rv.Interface().(protoEnum); ok && !c.UseEnumNumbers {
		return json.Marshal(enum.String())
	}
	return json.Marshal(rv.Interface())
}

// Unmarshal unmarshals JSON "data" into "v"
func (c *Codec) Unmarshal(data []byte, v any) error {
	return unmarshalJSONPb(data, c.UnmarshalOptions, v)
}

// NewDecoder returns a Decoder which reads JSON stream from "r".
func (c *Codec) NewDecoder(r io.Reader) codec.Decoder {
	d := json.NewDecoder(r)
	return DecoderWrapper{
		Decoder:          d,
		UnmarshalOptions: c.UnmarshalOptions,
	}
}

// DecoderWrapper is a wrapper around a *json.Decoder that adds
// support for protos to the Decode method.
type DecoderWrapper struct {
	*json.Decoder
	protojson.UnmarshalOptions
}

// Decode wraps the embedded decoder's Decode method to support
// protos using a jsonpb.Unmarshaler.
func (d DecoderWrapper) Decode(v any) error {
	return decodeJSONPb(d.Decoder, d.UnmarshalOptions, v)
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (c *Codec) NewEncoder(w io.Writer) codec.Encoder {
	return codec.EncoderFunc(func(v any) error {
		if err := c.marshalTo(w, v); err != nil {
			return err
		}
		// mimic json.Encoder by adding a newline (makes output
		// easier to read when it contains multiple encoded items)
		_, err := w.Write(c.Delimiter())
		return err
	})
}

func unmarshalJSONPb(data []byte, unmarshaler protojson.UnmarshalOptions, v any) error {
	d := json.NewDecoder(bytes.NewReader(data))
	return decodeJSONPb(d, unmarshaler, v)
}

func decodeJSONPb(d *json.Decoder, unmarshaler protojson.UnmarshalOptions, v any) error {
	p, ok := v.(proto.Message)
	if !ok {
		return decodeNonProtoField(d, unmarshaler, v)
	}

	// Decode into bytes for marshalling
	var b json.RawMessage
	err := d.Decode(&b)
	if err != nil {
		return err
	}

	return unmarshaler.Unmarshal([]byte(b), p)
}

func decodeNonProtoField(d *json.Decoder, unmarshaler protojson.UnmarshalOptions, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("%T is not a pointer", v)
	}
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		if rv.Type().ConvertibleTo(typeProtoMessage) {
			// Decode into bytes for marshalling
			var b json.RawMessage
			err := d.Decode(&b)
			if err != nil {
				return err
			}

			return unmarshaler.Unmarshal([]byte(b), rv.Interface().(proto.Message))
		}
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Map {
		if rv.IsNil() {
			rv.Set(reflect.MakeMap(rv.Type()))
		}
		conv, ok := convFromType[rv.Type().Key().Kind()]
		if !ok {
			return fmt.Errorf("unsupported type of map field key: %v", rv.Type().Key())
		}

		m := make(map[string]*json.RawMessage)
		if err := d.Decode(&m); err != nil {
			return err
		}
		for k, v := range m {
			result := conv.Call([]reflect.Value{reflect.ValueOf(k)})
			if err := result[1].Interface(); err != nil {
				return err.(error)
			}
			bk := result[0]
			bv := reflect.New(rv.Type().Elem())
			if v == nil {
				null := json.RawMessage("null")
				v = &null
			}
			if err := unmarshalJSONPb([]byte(*v), unmarshaler, bv.Interface()); err != nil {
				return err
			}
			rv.SetMapIndex(bk, bv.Elem())
		}
		return nil
	}
	if rv.Kind() == reflect.Slice {
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			var sl []byte
			if err := d.Decode(&sl); err != nil {
				return err
			}
			if sl != nil {
				rv.SetBytes(sl)
			}
			return nil
		}

		var sl []json.RawMessage
		if err := d.Decode(&sl); err != nil {
			return err
		}
		if sl != nil {
			rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
		}
		for _, item := range sl {
			bv := reflect.New(rv.Type().Elem())
			if err := unmarshalJSONPb([]byte(item), unmarshaler, bv.Interface()); err != nil {
				return err
			}
			rv.Set(reflect.Append(rv, bv.Elem()))
		}
		return nil
	}
	if _, ok := rv.Interface().(protoEnum); ok {
		var repr any
		if err := d.Decode(&repr); err != nil {
			return err
		}
		switch v := repr.(type) {
		case string:
			// TODO(yugui) Should use proto.StructProperties?
			return fmt.Errorf("unmarshaling of symbolic enum %q not supported: %T", repr, rv.Interface())
		case float64:
			rv.Set(reflect.ValueOf(int32(v)).Convert(rv.Type()))
			return nil
		default:
			return fmt.Errorf("cannot assign %#v into Go type %T", repr, rv.Interface())
		}
	}
	return d.Decode(v)
}

type protoEnum interface {
	fmt.Stringer
	EnumDescriptor() ([]byte, []int)
}

var typeProtoEnum = reflect.TypeOf((*protoEnum)(nil)).Elem()

var typeProtoMessage = reflect.TypeOf((*proto.Message)(nil)).Elem()

// Delimiter for newline encoded JSON streams.
func (c *Codec) Delimiter() []byte {
	return []byte("\n")
}

var (
	convFromType = map[reflect.Kind]reflect.Value{
		reflect.String:  reflect.ValueOf(codec.String),
		reflect.Bool:    reflect.ValueOf(codec.Bool),
		reflect.Float64: reflect.ValueOf(codec.Float64),
		reflect.Float32: reflect.ValueOf(codec.Float32),
		reflect.Int64:   reflect.ValueOf(codec.Int64),
		reflect.Int32:   reflect.ValueOf(codec.Int32),
		reflect.Uint64:  reflect.ValueOf(codec.Uint64),
		reflect.Uint32:  reflect.ValueOf(codec.Uint32),
		reflect.Slice:   reflect.ValueOf(codec.Bytes),
	}
)
