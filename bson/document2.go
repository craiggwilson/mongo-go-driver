// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bson

import (
	"bytes"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson/bsoncore"
	"github.com/mongodb/mongo-go-driver/bson/bsontype"
)

// Doc2 is a type safe, concise BSON document representation.
type Doc2 []Elem2

// ReadDoc2 will create a Document using the provided slice of bytes. If the
// slice of bytes is not a valid BSON document, this method will return an error.
func ReadDoc2(b []byte) (Doc2, error) {
	doc := make(Doc2, 0)
	err := doc.UnmarshalBSON(b)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// Copy makes a shallow copy of this document.
func (d Doc2) Copy() Doc2 {
	d2 := make(Doc2, len(d))
	copy(d2, d)
	return d2
}

// Append adds an element to the end of the document, creating it from the key and value provided.
func (d Doc2) Append(key string, val Val2) Doc2 {
	return append(d, Elem2{Key: key, Value: val})
}

// Prepend adds an element to the beginning of the document, creating it from the key and value provided.
func (d Doc2) Prepend(key string, val Val2) Doc2 {
	// TODO: should we just modify d itself instead of doing an alloc here?
	return append(Doc2{{Key: key, Value: val}}, d...)
}

// Set replaces an element of a document. If an element with a matching key is
// found, the element will be replaced with the one provided. If the document
// does not have an element with that key, the element is appended to the
// document instead.
func (d Doc2) Set(key string, val Val2) Doc2 {
	idx := d.indexOf(key)
	if idx == -1 {
		return append(d, Elem2{Key: key, Value: val})
	}
	d[idx] = Elem2{Key: key, Value: val}
	return d
}

func (d Doc2) indexOf(key string) int {
	for i, e := range d {
		if e.Key == key {
			return i
		}
	}
	return -1
}

// Delete removes the element with key if it exists and returns the updated Doc.
func (d Doc2) Delete(key string) Doc2 {
	idx := d.indexOf(key)
	if idx == -1 {
		return d
	}
	return append(d[:idx], d[idx+1:]...)
}

// Lookup searches the document and potentially subdocuments or arrays for the
// provided key. Each key provided to this method represents a layer of depth.
//
// This method will return an empty Value if they key does not exist. To know if they key actually
// exists, use LookupErr.
func (d Doc2) Lookup(key ...string) Val2 {
	val, _ := d.LookupErr(key...)
	return val
}

// LookupErr searches the document and potentially subdocuments or arrays for the
// provided key. Each key provided to this method represents a layer of depth.
func (d Doc2) LookupErr(key ...string) (Val2, error) {
	elem, err := d.LookupElementErr(key...)
	return elem.Value, err
}

// LookupElement searches the document and potentially subdocuments or arrays for the
// provided key. Each key provided to this method represents a layer of depth.
//
// This method will return an empty Element if they key does not exist. To know if they key actually
// exists, use LookupElementErr.
func (d Doc2) LookupElement(key ...string) Elem2 {
	elem, _ := d.LookupElementErr(key...)
	return elem
}

// LookupElementErr searches the document and potentially subdocuments for the
// provided key. Each key provided to this method represents a layer of depth.
func (d Doc2) LookupElementErr(key ...string) (Elem2, error) {
	// KeyNotFound operates by being created where the error happens and then the depth is
	// incremented by 1 as each function unwinds. Whenever this function returns, it also assigns
	// the Key slice to the key slice it has. This ensures that the proper depth is identified and
	// the proper keys.
	if len(key) == 0 {
		return Elem2{}, KeyNotFound{Key: key}
	}

	var elem Elem2
	var err error
	idx := d.indexOf(key[0])
	if idx == -1 {
		return Elem2{}, KeyNotFound{Key: key}
	}

	elem = d[idx]
	if len(key) == 1 {
		return elem, nil
	}

	switch elem.Value.Type() {
	case bsontype.EmbeddedDocument:
		switch tt := elem.Value.(type) {
		case Doc2:
			elem, err = tt.LookupElementErr(key[1:]...)
		}
	default:
		return Elem2{}, KeyNotFound{Type: elem.Value.Type()}
	}
	switch tt := err.(type) {
	case KeyNotFound:
		tt.Depth++
		tt.Key = key
		return Elem2{}, tt
	case nil:
		return elem, nil
	default:
		return Elem2{}, err // We can't actually hit this.
	}
}

// MarshalBSONValue implements the bsoncodec.ValueMarshaler interface.
//
// This method will never return an error.
func (d Doc2) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if d == nil {
		// TODO: Should we do this?
		return bsontype.Null, nil, nil
	}
	data, _ := d.MarshalBSON()
	return bsontype.EmbeddedDocument, data, nil
}

// MarshalBSON implements the Marshaler interface.
//
// This method will never return an error.
func (d Doc2) MarshalBSON() ([]byte, error) { return d.AppendMarshalBSON(nil) }

// AppendMarshalBSON marshals Doc to BSON bytes, appending to dst.
//
// This method will never return an error.
func (d Doc2) AppendMarshalBSON(dst []byte) ([]byte, error) {
	idx, dst := bsoncore.ReserveLength(dst)
	for _, elem := range d {
		t, data, _ := elem.Value.MarshalBSONValue() // Value.MarshalBSONValue never returns an error.
		dst = append(dst, byte(t))
		dst = append(dst, elem.Key...)
		dst = append(dst, 0x00)
		dst = append(dst, data...)
	}
	dst = append(dst, 0x00)
	dst = bsoncore.UpdateLength(dst, idx, int32(len(dst[idx:])))
	return dst, nil
}

// UnmarshalBSON implements the Unmarshaler interface.
func (d *Doc2) UnmarshalBSON(b []byte) error {
	if d == nil {
		return ErrNilDocument
	}

	if err := Raw(b).Validate(); err != nil {
		return err
	}

	elems, err := Raw(b).Elements()
	if err != nil {
		return err
	}
	var val Val2
	for _, elem := range elems {
		rawv := elem.Value()
		switch rawv.Type {
		case bsontype.EmbeddedDocument:
			val, _ = ReadDoc2(rawv.Value)
		default:
			v := Val{}
			v.UnmarshalBSONValue(rawv.Type, rawv.Value)
			val = v
		}
		if err != nil {
			return err
		}
		*d = d.Append(elem.Key(), val)
	}
	return nil
}

// Equal compares this document to another, returning true if they are equal.
func (d Doc2) Equal(v2 Val2) bool {
	return false
}

// String implements the fmt.Stringer interface.
func (d Doc2) String() string {
	var buf bytes.Buffer
	buf.Write([]byte("bson.Document{"))
	for idx, elem := range d {
		if idx > 0 {
			buf.Write([]byte(", "))
		}
		fmt.Fprintf(&buf, "%v", elem)
	}
	buf.WriteByte('}')

	return buf.String()
}

func (d Doc2) Type() bsontype.Type {
	return bsontype.EmbeddedDocument
}

func (Doc2) idoc() {}
