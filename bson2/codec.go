package bson2

import "reflect"

var globalRegistry = NewCodecRegistry()

func NewCodecRegistry() *CodecRegistry {
	codecs := map[reflect.Type]Codec{
		reflect.TypeOf(new(Document)): &DocumentCodec{},
		reflect.TypeOf(new(D)):        &DCodec{},
		reflect.TypeOf(new(M)):        &MCodec{},
		reflect.TypeOf(new(RawD)):     &RawDCodec{},
		reflect.TypeOf(new(Raw)):      &RawCodec{},
		reflect.TypeOf(new(int32)):    &Int32Codec{},
		reflect.TypeOf(new(string)):   &StringCodec{},
	}

	return &CodecRegistry{
		codecs: codecs,
	}
}

type CodecRegistry struct {
	codecs map[reflect.Type]Codec

	mapCodec    MapCodec
	structCodec StructCodec
}

func (cr *CodecRegistry) Lookup(t reflect.Type) (Codec, bool) {
	if t.Kind() != reflect.Ptr {
		t = reflect.PtrTo(t)
	}

	codec, ok := cr.codecs[t]
	if ok {
		return codec, true
	}

	// 1) See if it's a generic map type
	if t.Elem().Kind() == reflect.Map {
		return &cr.mapCodec, true
	}

	// 2) fallback to generic struct decoder
	if t.Elem().Kind() == reflect.Struct {
		return &cr.structCodec, true
	}
	return nil, false
}

func (cr *CodecRegistry) Register(t reflect.Type, codec Codec) {
	cr.codecs[t] = codec
}

type Codec interface {
	Decoder
	Encoder
}

type Decoder interface {
	Decode(reg *CodecRegistry, vr ValueReader, v interface{}) error
}

type Encoder interface {
	Encode(reg *CodecRegistry, vw ValueWriter, v interface{}) error
}
