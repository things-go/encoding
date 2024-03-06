package toml

import (
	"io"

	"github.com/pelletier/go-toml/v2"

	"github.com/things-go/encoding/codec"
)

// Codec is a Codec implementation with yaml.
type Codec struct{}

// ContentType always Returns "application/yaml; charset=utf-8".
func (*Codec) ContentType(_ any) string {
	return "application/toml; charset=utf-8"
}
func (*Codec) Marshal(v any) ([]byte, error) {
	return toml.Marshal(v)
}
func (*Codec) Unmarshal(data []byte, v any) error {
	return toml.Unmarshal(data, v)
}
func (*Codec) NewDecoder(r io.Reader) codec.Decoder {
	return toml.NewDecoder(r)
}
func (*Codec) NewEncoder(w io.Writer) codec.Encoder {
	return toml.NewEncoder(w)
}
