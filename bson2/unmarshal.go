package bson2

import (
	"fmt"
	"io"
	"reflect"
)

func Unmarshal(r io.Reader, v interface{}) error {
	return UnmarshalRegistry(r, v, globalRegistry)
}

func UnmarshalRegistry(r io.Reader, v interface{}, reg *CodecRegistry) error {
	t := reflect.TypeOf(v)
	codec, ok := reg.Lookup(t)
	if !ok {
		return fmt.Errorf("could not find codec for type %v", t)
	}

	vr, err := NewValueReader(r, TypeDocument)
	if err != nil {
		return err
	}

	_, err = codec.Decode(reg, vr, v)
	return err
}
