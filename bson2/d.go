package bson2

import "fmt"

type D []DocElem

type DocElem struct {
	Name  string
	Value interface{}
}

type DCodec struct{}

func (c *DCodec) Decode(reg *CodecRegistry, vr ValueReader, v interface{}) (interface{}, error) {
	var target *D
	var ok bool
	if v != nil {
		if target, ok = v.(*D); !ok {
			return nil, fmt.Errorf("%T can only be used to decode *bson.D", c)
		}
	} else {
		target = &D{}
	}

	elems := []DocElem(*target)

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

		elems = append(elems, DocElem{name, value})
	}

	*target = D(elems)
	return target, nil
}

func (c *DCodec) decodeValue(reg *CodecRegistry, vr ValueReader) (interface{}, error) {
	switch vr.Type() {
	case TypeBoolean:
		return vr.ReadBoolean()
	case TypeDocument:
		value, err := c.Decode(reg, vr, nil)
		if err != nil {
			return nil, err
		}

		return *value.(*D), nil
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
