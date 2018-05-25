package bson2

import "fmt"

type M map[string]interface{}

type MCodec struct{}

func (c *MCodec) Decode(reg *CodecRegistry, vr ValueReader, v interface{}) (interface{}, error) {
	var target *M
	var ok bool
	if v != nil {
		if target, ok = v.(*M); !ok {
			return nil, fmt.Errorf("%T can only be used to decode *bson.D", c)
		}
	} else {
		target = &M{}
	}

	m := map[string]interface{}(*target)

	doc, err := vr.ReadDocument()
	if err != nil {
		return nil, err
	}

	for {
		name, evr, err := doc.ReadElement()
		if err == EOD {
			break
		}
		if err != nil {
			return nil, err
		}

		value, err := c.decodeValue(reg, evr)
		if err != nil {
			return nil, err
		}

		m[name] = value
	}

	*target = m
	return target, nil
}

func (c *MCodec) decodeValue(reg *CodecRegistry, vr ValueReader) (interface{}, error) {
	switch vr.Type() {
	case TypeBoolean:
		return vr.ReadBoolean()
	case TypeDocument:
		value, err := c.Decode(reg, vr, nil)
		if err != nil {
			return nil, err
		}

		return *value.(*M), nil
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
