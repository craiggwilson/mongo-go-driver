package bson2

import (
	"fmt"
	"reflect"
)

type MapCodec struct{}

func (c *MapCodec) Decode(reg *CodecRegistry, vr ValueReader, v interface{}) error {
	mPtr := reflect.ValueOf(v)
	if mPtr.Kind() != reflect.Ptr {
		return fmt.Errorf("%T can only process pointers to maps, but got a %T", c, v)
	}

	m := mPtr.Elem()
	mType := m.Type()
	keyType := mType.Key()
	valueType := mType.Elem()

	dr, err := vr.ReadDocument()
	if err != nil {
		return err
	}
	for {
		name, vr, err := dr.ReadElement()
		if err == EOD {
			break
		}
		if err != nil {
			return err
		}

		keyValue, err := c.nameToKey(name, keyType)
		if err != nil {
			return err
		}

		value := m.MapIndex(keyValue)
		if value == zeroVal {
			value, err = c.decodeValue(reg, vr, valueType, mType)
			if err != nil {
				return err
			}
			m.SetMapIndex(keyValue, value)
		} else {
			panic("not supported")
		}
	}

	return nil
}

func (c *MapCodec) nameToKey(name string, keyType reflect.Type) (reflect.Value, error) {
	switch keyType {
	case tString:
		return reflect.ValueOf(name), nil
	case tIface:
		return reflect.ValueOf(name), nil
	default:
		return zeroVal, fmt.Errorf("cannot decode into a map with a key of type %s", keyType)
	}
}

func (c *MapCodec) decodeValue(reg *CodecRegistry, vr ValueReader, valueType reflect.Type, mType reflect.Type) (reflect.Value, error) {
	switch valueType {
	case tIface:
		return decodeIface(reg, vr, mType)
	default:
		return zeroVal, fmt.Errorf("unsupported AHAHAHAHA")
	}
}
