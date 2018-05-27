package bson2

import (
	"fmt"
)

type M map[string]interface{}

type MCodec struct{}

func (c *MCodec) Decode(reg *CodecRegistry, vr ValueReader, v interface{}) error {
	var target *M
	var ok bool
	if target, ok = v.(*M); !ok {
		return fmt.Errorf("%T can only be used to decode *bson2.M", c)
	}

	m := map[string]interface{}(*target)

	doc, err := vr.ReadDocument()
	if err != nil {
		return err
	}

	for {
		name, vr, err := doc.ReadElement()
		if err == EOD {
			break
		}
		if err != nil {
			return err
		}

		// NOTE: we could call decodeIface here and share some code, but it seems that it has a
		// significant impact on performance, presumably due to it's use of reflection, so...
		value, err := c.decodeValue(reg, vr)
		if err != nil {
			return err
		}

		m[name] = value
	}

	*target = m
	return nil
}

func (c *MCodec) decodeValue(reg *CodecRegistry, vr ValueReader) (interface{}, error) {
	switch vr.Type() {
	case TypeBoolean:
		return vr.ReadBoolean()
	case TypeDocument:
		value := M{}
		if err := c.Decode(reg, vr, &value); err != nil {
			return nil, err
		}

		return value, nil
	case TypeInt32:
		return vr.ReadInt32()
	case TypeInt64:
		return vr.ReadInt64()
	case TypeString:
		return vr.ReadString()
	default:
		return nil, fmt.Errorf("unsuppored bson type")
	}
}
