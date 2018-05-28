package bson2

import (
	"bytes"
	"fmt"
	"reflect"
)

func Marshal(v interface{}) ([]byte, error) {
	return MarshalRegistry(v, globalRegistry)
}

func MarshalRegistry(v interface{}, reg *CodecRegistry) ([]byte, error) {
	t := reflect.TypeOf(v)
	codec, ok := reg.Lookup(t)
	if !ok {
		return nil, fmt.Errorf("could not find codec for type %v", t)
	}

	vw := NewValueWriter()
	if err := codec.Encode(reg, vw, v); err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	_, err := vw.WriteTo(&buffer)
	return buffer.Bytes(), err
}
