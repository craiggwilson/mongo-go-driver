package bson2

import "fmt"

type StringCodec struct{}

func (c *StringCodec) Decode(reg *CodecRegistry, vr ValueReader, v interface{}) error {
	var target *string
	var ok bool
	if target, ok = v.(*string); !ok {
		return fmt.Errorf("%T can only be used to decode *string", c)
	}

	var err error
	switch vr.Type() {
	case TypeString:
		*target, err = vr.ReadString()
		return err
	default:
		return fmt.Errorf("cannot decode %v into a string", vr.Type())
	}
}
