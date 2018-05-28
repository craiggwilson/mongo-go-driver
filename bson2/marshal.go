package bson2

import (
	"fmt"
	"io"
	"reflect"
)

func Marshal(v interface{}, out io.Writer) error {
	return MarshalRegistry(v, out, globalRegistry)
}

func MarshalRegistry(v interface{}, out io.Writer, reg *CodecRegistry) error {
	t := reflect.TypeOf(v)
	codec, ok := reg.Lookup(t)
	if !ok {
		return fmt.Errorf("could not find codec for type %v", t)
	}

	vw := NewValueWriter()
	if err := codec.Encode(reg, vw, v); err != nil {
		return err
	}

	_, err := vw.WriteTo(out)
	return err
}
