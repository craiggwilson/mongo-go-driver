package bson2

import (
	"fmt"
	"reflect"
)

type RawD []RawDocElem

type RawDocElem struct {
	Name  string
	Value Raw
}

type Raw struct {
	Kind byte
	Data []byte

	reg *CodecRegistry
}

func (raw Raw) Unmarshal(v interface{}) error {
	t := reflect.TypeOf(v)
	codec, ok := raw.reg.Lookup(t)
	if !ok {
		return fmt.Errorf("could not find codec for type %v", t)
	}

	vr, err := NewValueReader(raw.Data, Type(raw.Kind))
	if err != nil {
		return err
	}

	return codec.Decode(raw.reg, vr, v)
}

type RawDCodec struct{}

func (c *RawDCodec) Decode(reg *CodecRegistry, vr ValueReader, v interface{}) error {
	var target *RawD
	var ok bool
	if target, ok = v.(*RawD); !ok {
		return fmt.Errorf("%T can only be used to decode *bson2.RawD", c)
	}

	elems := []RawDocElem(*target)

	doc, err := vr.ReadDocument()
	if err != nil {
		return err
	}

	for {
		name, evr, err := doc.ReadElement()
		if err == EOD {
			break
		}
		if err != nil {
			return err
		}

		kind := byte(evr.Type())

		bytes := make([]byte, evr.Size())
		err = evr.ReadBytes(bytes)
		if err != nil {
			return err
		}

		elems = append(elems, RawDocElem{Name: name, Value: Raw{Kind: kind, Data: bytes, reg: reg}})
	}

	*target = RawD(elems)
	return nil
}

func (c *RawDCodec) Encode(reg *CodecRegistry, vw ValueWriter, v interface{}) error {
	return fmt.Errorf("not supported")
}

type RawCodec struct{}

func (c *RawCodec) Decode(reg *CodecRegistry, vr ValueReader, v interface{}) error {
	var target *Raw
	var ok bool
	if target, ok = v.(*Raw); !ok {
		return fmt.Errorf("%T can only be used to decode *bson2.Raw", c)
	}

	target.Kind = byte(vr.Type())
	target.Data = make([]byte, vr.Size())
	vr.ReadBytes(target.Data)
	target.reg = reg
	return nil
}

func (c *RawCodec) Encode(reg *CodecRegistry, vw ValueWriter, v interface{}) error {
	return fmt.Errorf("not supported")
}
