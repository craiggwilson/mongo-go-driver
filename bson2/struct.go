package bson2

import (
	"fmt"
	"reflect"
	"strings"
)

type StructCodec struct{}

var zeroVal reflect.Value

func (c *StructCodec) Decode(reg *CodecRegistry, vr ValueReader, v interface{}) error {
	valuePtr := reflect.ValueOf(v)
	if valuePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("%T can only process pointers to structs, but got a %T", c, v)
	}

	value := valuePtr.Elem()

	valueType := value.Type()

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

		field := value.FieldByNameFunc(func(field string) bool {
			return c.matchesField(name, field, valueType)
		})
		if field == zeroVal {
			// TODO: the default should not be ignore, although that is currently how mgo works.
			// We should return an error unless they have:
			// 1) specified to ignore extra elements
			// 2) included an extra elements field (map[string]interface{}, bson.D, bson.Document, etc...)
			continue
		}

		fieldPtr := field
		if fieldPtr.Kind() != reflect.Ptr {
			if !field.CanAddr() {
				return fmt.Errorf("cannot decode element '%s' into field %v; it is not addressable", name, field)
			}
			fieldPtr = field.Addr()
		} else if fieldPtr.IsNil() {
			if !fieldPtr.CanSet() {
				return fmt.Errorf("cannot decode element '%s' into field %v; it is not settable", name, field)
			}
			fieldPtr.Set(reflect.New(fieldPtr.Type().Elem()))
		}

		fieldPtrType := fieldPtr.Type()
		fieldPtrCodec, ok := reg.Lookup(fieldPtrType)
		if !ok {
			return fmt.Errorf("unable to find codec for type %v for field '%s'", name, fieldPtrType)
		}

		if err = fieldPtrCodec.Decode(reg, vr, fieldPtr.Interface()); err != nil {
			return err
		}
	}

	return nil
}

func (c *StructCodec) matchesField(key string, field string, sType reflect.Type) bool {
	sField, found := sType.FieldByName(field)
	if !found {
		return false
	}

	tag, ok := sField.Tag.Lookup("bson")
	if !ok {
		// Get the full tag string
		tag = string(sField.Tag)

		if len(sField.Tag) == 0 || strings.ContainsRune(tag, ':') {
			return strings.ToLower(key) == strings.ToLower(field)
		}
	}

	var fieldKey string
	i := strings.IndexRune(tag, ',')
	if i == -1 {
		fieldKey = tag
	} else {
		fieldKey = tag[:i]
	}

	return fieldKey == key
}
