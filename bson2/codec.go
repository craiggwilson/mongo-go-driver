package bson2

import "reflect"

var globalRegistry = NewCodecRegistry()

func NewCodecRegistry() *CodecRegistry {
	codecs := map[reflect.Type]Codec{
		reflect.TypeOf(&D{}):    &DCodec{},
		reflect.TypeOf(&M{}):    &MCodec{},
		reflect.TypeOf(&RawD{}): &RawDCodec{},
		reflect.TypeOf(&Raw{}):  &RawCodec{},
	}

	return &CodecRegistry{
		codecs: codecs,
	}
}

type CodecRegistry struct {
	codecs map[reflect.Type]Codec
}

func (cr *CodecRegistry) Lookup(t reflect.Type) (Codec, bool) {
	codec, ok := cr.codecs[t]
	return codec, ok
}

func (cr *CodecRegistry) Register(t reflect.Type, codec Codec) {
	cr.codecs[t] = codec
}

type Codec interface {
	Decoder
}

type Decoder interface {
	Decode(reg *CodecRegistry, vr ValueReader, v interface{}) (interface{}, error)
}
