package bson2

import (
	"encoding/binary"
	"fmt"
)

type Document struct {
	Elements []*Element
}

type Element struct {
	Name  string
	Value *Value
}

type Value struct {
	t     Type
	bytes []byte
}

func (v *Value) Bytes() []byte {
	return v.bytes
}

func (v *Value) Type() Type {
	return v.t
}

func (v *Value) Boolean() (bool, error) {
	if v.t != TypeBoolean {
		return false, newErrValueType(v.t, TypeBoolean)
	}

	if v.bytes[0] == 0 {
		return false, nil
	} else if v.bytes[0] == 1 {
		return true, nil
	}

	return false, errInvalidValue(fmt.Sprintf("invalid byte for boolean, %s", v.bytes[0]))
}

func (v *Value) Int32() (int32, error) {
	if v.t != TypeInt32 {
		return 0, newErrValueType(v.t, TypeInt32)
	}
	return int32(binary.LittleEndian.Uint32(v.bytes)), nil
}

func (v *Value) Int64() (int64, error) {
	if v.t != TypeInt64 {
		return 0, newErrValueType(v.t, TypeInt64)
	}
	return int64(binary.LittleEndian.Uint64(v.bytes)), nil
}

func (v *Value) StringValue() (string, error) {
	if v.t != TypeString {
		return "", newErrValueType(v.t, TypeString)
	}
	return string(v.bytes), nil
}

type DocumentCodec struct{}

func (c *DocumentCodec) Decode(r ValueReader, value interface{}) error {

	document, ok := value.(*Document)
	if !ok {
		return fmt.Errorf("value is a %T, but expected a *bson.Document", value)
	}

	dr, err := r.ReadDocument()
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

		bytes := make([]byte, vr.Size())
		err = vr.ReadBytes(bytes)
		if err != nil {
			return err
		}

		document.Elements = append(document.Elements, &Element{
			Name: name,
			Value: &Value{
				t:     vr.Type(),
				bytes: bytes,
			},
		})
	}

	return nil
}
