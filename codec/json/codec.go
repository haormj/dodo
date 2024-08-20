package json

import (
	"encoding/json"

	"github.com/haormj/dodo/codec"
)

// Codec implement by json
type Codec struct{}

func NewCodec() codec.Codec {
	return Codec{}
}

func (Codec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (Codec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (Codec) String() string {
	return "json"
}
