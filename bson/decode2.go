package bson

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"reflect"
)

func NewCodecRegistry() *CodecRegistry {
	codecs := map[reflect.Type]Codec2{
		reflect.TypeOf(""):       &StringCodec{},
		reflect.TypeOf(int32(0)): &Int32Codec{},
	}

	return &CodecRegistry{
		codecs: codecs,
	}
}

type CodecRegistry struct {
	codecs map[reflect.Type]Codec2
}

func (cr *CodecRegistry) Lookup(t reflect.Type) Codec2 {
	return cr.codecs[t]
}

func (cr *CodecRegistry) Registry(t reflect.Type, codec Codec2) {
	cr.codecs[t] = codec
}

type Decoder2 interface {
	Decode(*Value, interface{}) (interface{}, error)
}

type Codec2 interface {
	Decoder2
}

func NewCodecRegistryCodec(registry *CodecRegistry) *CodecRegistryCodec {
	return &CodecRegistryCodec{registry: registry}
}

type CodecRegistryCodec struct {
	registry *CodecRegistry
}

func (cr *CodecRegistryCodec) Decode(r io.Reader, v interface{}) (interface{}, error) {
	d := &decoder2{
		registry:      cr.registry,
		r:             newPeekLengthReader(r),
		containerType: reflect.TypeOf(make(map[string]interface{})),
	}
	return v, d.decode(v)
}

type decoder2 struct {
	registry *CodecRegistry
	r        *peekLengthReader

	containerType reflect.Type
}

func (d *decoder2) createEmpty(t reflect.Type) reflect.Value {
	switch t.Kind() {
	case reflect.Map:
		return reflect.MakeMap(t)
	case reflect.Struct:
		return reflect.New(t).Elem()
	case reflect.Ptr:
		empty := d.createEmpty(t.Elem())
		value := reflect.New(empty.Type())
		value.Elem().Set(empty)
		return value
	default:
		panic(fmt.Sprintf("create empty not supported for %v", t))
	}
}

func (d *decoder2) decode(v interface{}) error {
	switch t := v.(type) {
	case []byte:
		length, err := d.r.peekLength()
		if err != nil {
			return err
		}

		if len(t) < int(length) {
			return NewErrTooSmall()
		}

		_, err = io.ReadFull(d.r, t)
		if err != nil {
			return err
		}

		_, err = Reader(t).Validate()
		return err
	default:
		target := reflect.ValueOf(v)
		return d.reflectDecode(target)
	}
}

func (d *decoder2) reflectDecode(target reflect.Value) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s", e)
		}
	}()

	switch target.Kind() {
	case reflect.Map:
		if target == zeroVal {
			target = d.createEmpty(d.containerType)
		}
		return d.decodeMap(target)
	case reflect.Struct:
		return d.decodeStruct(target)
	case reflect.Ptr:
		return d.reflectDecode(target.Elem())
	default:
		return fmt.Errorf("unsupported target %v", target.Kind())
	}
}

func (d *decoder2) decodeMap(target reflect.Value) (err error) {
	oldContainerType := d.containerType
	d.containerType = target.Type()

	r, err := NewFromIOReader(d.r)
	if err != nil {
		return err
	}

	iter, err := r.Iterator()
	if err != nil {
		return err
	}

	for iter.Next() {
		elem := iter.Element()

		key := reflect.ValueOf(elem.Key())
		value := target.MapIndex(key)
		var valueType reflect.Type
		if value == zeroVal {
			valueType = d.containerType.Elem()
		} else {
			valueType = value.Type()
		}

		value, err = d.reflectDecodeValue(elem.value, value, valueType)
		if err != nil {
			return err
		}

		target.SetMapIndex(key, value)
	}

	d.containerType = oldContainerType

	return iter.Err()
}

func (d *decoder2) decodeStruct(target reflect.Value) error {
	oldContainerType := d.containerType
	d.containerType = reflect.TypeOf(make(map[string]interface{}))

	r, err := NewFromIOReader(d.r)
	if err != nil {
		return err
	}

	iter, err := r.Iterator()
	if err != nil {
		return err
	}

	for iter.Next() {
		elem := iter.Element()

		field := target.FieldByNameFunc(func(field string) bool {
			return matchesField(elem.Key(), field, target.Type())
		})
		if field == zeroVal {
			continue
		}

		value := field
		if field.Kind() == reflect.Struct {
			value = value.Addr()
		}

		value, err := d.reflectDecodeValue(elem.value, value, field.Type())
		if err != nil {
			return err
		}

		if value != zeroVal {
			if field.Kind() == reflect.Struct {
				value = value.Elem()
			}
			field.Set(value)
		}
	}

	d.containerType = oldContainerType

	return iter.Err()
}

func (d *decoder2) reflectDecodeValue(bsonValue *Value, value reflect.Value, valueType reflect.Type) (reflect.Value, error) {
	if codec, ok := d.registry.codecs[valueType]; ok {
		i, err := codec.Decode(bsonValue, value)
		return reflect.ValueOf(i), err
	}

	// this is fallback when we don't know what type we are decoding to. It will usually occur
	// when decoding into a map.
	switch bsonValue.Type() {
	case TypeString:
		s := bsonValue.StringValue()
		switch valueType {
		case tString, tEmpty:
			return reflect.ValueOf(s), nil
		}
	case TypeInt32:
		i32 := bsonValue.Int32()
		switch valueType {
		case tString, tEmpty:
			return reflect.ValueOf(i32), nil
		}
	case TypeEmbeddedDocument:
		r := bsonValue.ReaderDocument()
		newD := &decoder2{
			r:             newPeekLengthReader(bytes.NewBuffer(r)),
			registry:      d.registry,
			containerType: d.containerType,
		}

		if value == zeroVal || (value.Kind() == reflect.Ptr && value.IsNil()) {
			if valueType == tEmpty {
				valueType = d.containerType
			}

			value = d.createEmpty(valueType)
		}

		err := newD.decode(value.Interface())
		return value, err
	}

	return value, fmt.Errorf("unsupported combination, %s and %s", bsonValue.Type(), value)
}

type Int32Codec struct{}

func (c *Int32Codec) Decode(value *Value, _ interface{}) (interface{}, error) {
	switch value.Type() {
	case TypeInt32:
		return value.Int32(), nil
	case TypeInt64:
		i64 := value.Int64()
		if i64 > math.MaxInt32 {
			return nil, fmt.Errorf("overflow error")
		}
		return int32(i64), nil
	default:
		return nil, fmt.Errorf("Int32Codec cannot decode a %s", value.Type())
	}
}

type StringCodec struct{}

func (c *StringCodec) Decode(value *Value, _ interface{}) (interface{}, error) {
	switch value.Type() {
	case TypeString:
		return value.StringValue(), nil
	default:
		return nil, fmt.Errorf("StringCodec cannot decode a %s", value.Type())
	}
}
