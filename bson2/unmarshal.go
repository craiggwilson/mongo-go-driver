package bson2

import (
	"fmt"
	"reflect"
)

func Unmarshal(input []byte, v interface{}) error {
	return UnmarshalRegistry(input, v, globalRegistry)
}

func UnmarshalRegistry(input []byte, v interface{}, reg *CodecRegistry) error {
	t := reflect.TypeOf(v)
	codec, ok := reg.Lookup(t)
	if !ok {
		return fmt.Errorf("could not find codec for type %v", t)
	}

	vr, err := NewValueReader(input, TypeDocument)
	if err != nil {
		return err
	}

	_, err = codec.Decode(reg, vr, v)
	return err
}
