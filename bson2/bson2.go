package bson2

import (
	"fmt"
	"reflect"
)

var zeroVal reflect.Value

var tString = reflect.TypeOf("")
var tIface = reflect.TypeOf((*interface{})(nil)).Elem()

func decodeIface(reg *CodecRegistry, vr ValueReader, containerType reflect.Type) (reflect.Value, error) {
	switch vr.Type() {
	case TypeBoolean:
		b, err := vr.ReadBoolean()
		if err != nil {
			return zeroVal, err
		}
		return reflect.ValueOf(b), nil
	case TypeDocument:
		d, err := createEmptyContainer(containerType)
		if err != nil {
			return zeroVal, err
		}
		dPtr := d
		if dPtr.Kind() != reflect.Ptr {
			dPtr = reflect.New(containerType)
			dPtr.Elem().Set(d)
		}

		codec, ok := reg.Lookup(dPtr.Type())
		if !ok {
			return zeroVal, fmt.Errorf("could not find codec for type %s", dPtr.Type())
		}

		if err := codec.Decode(reg, vr, dPtr.Interface()); err != nil {
			return zeroVal, err
		}

		return d, nil

	case TypeInt32:
		i32, err := vr.ReadInt32()
		if err != nil {
			return zeroVal, err
		}
		return reflect.ValueOf(i32), nil
	case TypeInt64:
		i64, err := vr.ReadInt64()
		if err != nil {
			return zeroVal, err
		}
		return reflect.ValueOf(i64), err
	case TypeString:
		s, err := vr.ReadString()
		if err != nil {
			return zeroVal, err
		}
		return reflect.ValueOf(s), nil
	default:
		return zeroVal, fmt.Errorf("unsuppored bson type")
	}
}

func createEmptyContainer(containerType reflect.Type) (reflect.Value, error) {
	switch containerType.Kind() {
	case reflect.Ptr:
		container, err := createEmptyContainer(containerType.Elem())
		if err != nil {
			return zeroVal, err
		}

		containerPtr := reflect.New(container.Type())
		containerPtr.Elem().Set(container)

		return containerPtr, nil
	case reflect.Map:
		return reflect.MakeMap(containerType), nil
	default:
		return zeroVal, fmt.Errorf("unsupported container type %v", containerType)
	}
}
