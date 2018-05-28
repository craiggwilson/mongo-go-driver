package bson2

import (
	"fmt"
	"math"
)

type Int32Codec struct{}

func (c *Int32Codec) Decode(reg *CodecRegistry, vr ValueReader, v interface{}) error {
	var target *int32
	var ok bool
	if target, ok = v.(*int32); !ok {
		return fmt.Errorf("%T can only be used to decode *string", c)
	}

	var err error
	switch vr.Type() {
	case TypeInt32:
		*target, err = vr.ReadInt32()
		return err
	case TypeInt64:
		i64, err := vr.ReadInt64()
		if err != nil {
			return err
		}

		if i64 > math.MaxInt32 {
			return fmt.Errorf("overflow detected")
		}

		*target = int32(i64)
		return nil
	default:
		return fmt.Errorf("cannot decode %v into a string", vr.Type())
	}
}

func (c *Int32Codec) Encode(reg *CodecRegistry, vw ValueWriter, v interface{}) error {
	var value int32
	if valuePtr, ok := v.(*int32); ok {
		value = *valuePtr
	} else if value, ok = v.(int32); !ok {
		return fmt.Errorf("%T can only be used to encode int32 or *int32", c)
	}

	return vw.WriteInt32(value)
}

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

func (c *StringCodec) Encode(reg *CodecRegistry, vw ValueWriter, v interface{}) error {
	var value string
	if valuePtr, ok := v.(*string); ok {
		value = *valuePtr
	} else if value, ok = v.(string); !ok {
		return fmt.Errorf("%T can only be used to encode string or *string", c)
	}

	return vw.WriteString(value)
}
