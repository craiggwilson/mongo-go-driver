package bson2

import "fmt"

type D []DocElem

type DocElem struct {
	Name  string
	Value interface{}
}

type DCodec struct{}

func (c *DCodec) Decode(reg *CodecRegistry, vr ValueReader, v interface{}) error {
	var target *D
	var ok bool
	if target, ok = v.(*D); !ok {
		return fmt.Errorf("%T can only be used to decode *bson2.D", c)
	}

	elems := []DocElem(*target)

	dr, err := vr.ReadDocument()
	if err != nil {
		return err
	}

	for {
		name, evr, err := dr.ReadElement()
		if err == EOD {
			break
		}
		if err != nil {
			return err
		}

		value, err := c.decodeValue(reg, evr)
		if err != nil {
			return err
		}

		elems = append(elems, DocElem{name, value})
	}

	*target = D(elems)
	return nil
}

func (c *DCodec) decodeValue(reg *CodecRegistry, vr ValueReader) (interface{}, error) {
	switch vr.Type() {
	case TypeBoolean:
		return vr.ReadBoolean()
	case TypeDocument:
		value := D{}
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
