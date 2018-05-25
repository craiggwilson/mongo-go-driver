package bson2

import (
	"bytes"
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

	r := bytes.NewReader(raw.Data)
	vr, err := NewValueReaderFromIO(r, Type(raw.Kind))
	if err != nil {
		return err
	}

	_, err = codec.Decode(raw.reg, vr, v)
	return err
}

type RawDCodec struct{}

func (c *RawDCodec) Decode(reg *CodecRegistry, vr ValueReader, v interface{}) (interface{}, error) {
	var target *RawD
	var ok bool
	if v != nil {
		if target, ok = v.(*RawD); !ok {
			return nil, fmt.Errorf("%T can only be used to decode *bson.D", c)
		}
	} else {
		target = &RawD{}
	}

	elems := []RawDocElem(*target)

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

		kind := byte(evr.Type())

		bytes := make([]byte, evr.Size())
		err = evr.ReadBytes(bytes)
		if err != nil {
			return nil, err
		}

		elems = append(elems, RawDocElem{Name: name, Value: Raw{Kind: kind, Data: bytes, reg: reg}})
	}

	*target = RawD(elems)
	return target, nil
}