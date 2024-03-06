package xml

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

type Plain struct {
	V any
}

type NestedOrder struct {
	XMLName struct{} `xml:"result"`
	Field1  string   `xml:"parent>c"`
	Field2  string   `xml:"parent>b"`
	Field3  string   `xml:"parent>a"`
}

func TestCodec_ContentType(t *testing.T) {
	codec := Codec{}

	want := "application/xml; charset=utf-8"
	if got := codec.ContentType(struct{}{}); got != want {
		t.Errorf("m.ContentType(_) failed, got = %q; want %q; ", got, want)
	}
}

func TestCodec_Marshal(t *testing.T) {
	codec := Codec{}

	tests := []struct {
		Value     any
		ExpectXML string
	}{
		// Test value types
		{Value: &Plain{true}, ExpectXML: `<Plain><V>true</V></Plain>`},
		{Value: &Plain{false}, ExpectXML: `<Plain><V>false</V></Plain>`},
		{Value: &Plain{42}, ExpectXML: `<Plain><V>42</V></Plain>`},
		{
			Value: &NestedOrder{Field1: "C", Field2: "B", Field3: "A"},
			ExpectXML: `<result>` +
				`<parent>` +
				`<c>C</c>` +
				`<b>B</b>` +
				`<a>A</a>` +
				`</parent>` +
				`</result>`,
		},
	}
	for _, tt := range tests {
		data, err := codec.Marshal(tt.Value)
		if err != nil {
			t.Errorf("Marshal(%#v): %s", tt.Value, err)
		}
		if got, want := string(data), tt.ExpectXML; got != want {
			if strings.Contains(want, "\n") {
				t.Errorf("marshal(%#v):\nHAVE:\n%s\nWANT:\n%s", tt.Value, got, want)
			} else {
				t.Errorf("marshal(%#v):\nhave %#q\nwant %#q", tt.Value, got, want)
			}
		}
		data1 := &bytes.Buffer{}
		err = codec.NewEncoder(data1).Encode(tt.Value)
		if err != nil {
			t.Errorf("Encode(%#v): %s", tt.Value, err)
		}
		if got, want := data1.String(), tt.ExpectXML; got != want {
			if strings.Contains(want, "\n") {
				t.Errorf("Encode(%#v):\nHAVE:\n%s\nWANT:\n%s", tt.Value, got, want)
			} else {
				t.Errorf("Encode(%#v):\nhave %#q\nwant %#q", tt.Value, got, want)
			}
		}
	}
}

func TestCodec_Unmarshal(t *testing.T) {
	codec := Codec{}

	tests := []struct {
		want     any
		InputXML string
	}{
		{
			want: &NestedOrder{Field1: "C", Field2: "B", Field3: "A"},
			InputXML: `<result>` +
				`<parent>` +
				`<c>C</c>` +
				`<b>B</b>` +
				`<a>A</a>` +
				`</parent>` +
				`</result>`,
		},
	}

	for _, tt := range tests {
		vt := reflect.TypeOf(tt.want)
		dest := reflect.New(vt.Elem()).Interface()
		data := []byte(tt.InputXML)
		err := codec.Unmarshal(data, dest)
		if err != nil {
			t.Errorf("unmarshal(%#v, %#v): %s", tt.InputXML, dest, err)
		}
		if got, want := dest, tt.want; !reflect.DeepEqual(got, want) {
			t.Errorf("unmarshal(%q):\nhave %#v\nwant %#v", tt.InputXML, got, want)
		}

		dest1 := reflect.New(vt.Elem()).Interface()
		err = codec.NewDecoder(bytes.NewBuffer(data)).Decode(dest1)
		if err != nil {
			t.Errorf("Decode(%#v, %#v): %s", tt.InputXML, dest, err)
		}
		if got, want := dest1, tt.want; !reflect.DeepEqual(got, want) {
			t.Errorf("Decode(%q):\nhave %#v\nwant %#v", tt.InputXML, got, want)
		}
	}
}

func TestCodec_NilUnmarshal(t *testing.T) {
	codec := Codec{}
	tests := []struct {
		want     any
		InputXML string
	}{
		{
			want: &NestedOrder{Field1: "C", Field2: "B", Field3: "A"},
			InputXML: `<result>` +
				`<parent>` +
				`<c>C</c>` +
				`<b>B</b>` +
				`<a>A</a>` +
				`</parent>` +
				`</result>`,
		},
	}

	for _, tt := range tests {
		s := struct {
			A string `xml:"a"`
			B *NestedOrder
		}{A: "a"}
		data := []byte(tt.InputXML)
		err := codec.Unmarshal(data, &s.B)
		if err != nil {
			t.Errorf("unmarshal(%#v, %#v): %s", tt.InputXML, s.B, err)
		}
		if got, want := s.B, tt.want; !reflect.DeepEqual(got, want) {
			t.Errorf("unmarshal(%q):\nhave %#v\nwant %#v", tt.InputXML, got, want)
		}
	}
}
