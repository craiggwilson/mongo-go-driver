package bson2

import (
	"fmt"
	"reflect"
)

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

func (c *DCodec) Encode(reg *CodecRegistry, vw ValueWriter, v interface{}) error {
	var value D
	if valuePtr, ok := v.(*D); ok {
		value = *valuePtr
	} else if value, ok = v.(D); !ok {
		return fmt.Errorf("%T can only be used to encode bson2.D or *bson.D", c)
	}

	elems := []DocElem(value)

	dw, err := vw.WriteDocument()
	if err != nil {
		return err
	}

	for _, elem := range elems {
		vw, err := dw.WriteElement(elem.Name)
		if err != nil {
			return err
		}

		codec, ok := reg.Lookup(reflect.TypeOf(elem.Value))
		if !ok {
			return fmt.Errorf("could not find codec for %T", elem.Value)
		}

		if err = codec.Encode(reg, vw, elem.Value); err != nil {
			return err
		}
	}

	return dw.WriteEndDocument()
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
