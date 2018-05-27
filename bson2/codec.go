package bson2

import "reflect"

var globalRegistry = NewCodecRegistry()

func NewCodecRegistry() *CodecRegistry {
	codecs := map[reflect.Type]Codec{
		reflect.TypeOf(&Document{}): &DocumentCodec{},
		reflect.TypeOf(&D{}):        &DCodec{},
		reflect.TypeOf(&M{}):        &MCodec{},
		reflect.TypeOf(&RawD{}):     &RawDCodec{},
		reflect.TypeOf(&Raw{}):      &RawCodec{},
		reflect.TypeOf(new(string)): &StringCodec{},
	}

	return &CodecRegistry{
		codecs: codecs,
	}
}

type CodecRegistry struct {
	codecs map[reflect.Type]Codec

	structCodec StructCodec
}

func (cr *CodecRegistry) Lookup(t reflect.Type) (Codec, bool) {
	codec, ok := cr.codecs[t]
	if !ok {
		// 1) check to see if this type implements Unmarshaler. If so, invoke that method.
		// TODO
		// 2) fallback to generic struct decoder
		if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
			return &cr.structCodec, true
		}
	}
	return codec, ok
}

func (cr *CodecRegistry) Register(t reflect.Type, codec Codec) {
	cr.codecs[t] = codec
}

type Codec interface {
	Decoder
}

type Decoder interface {
	Decode(reg *CodecRegistry, vr ValueReader, v interface{}) error
}
