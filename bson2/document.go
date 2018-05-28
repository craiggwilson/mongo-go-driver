package bson2

import (
	"bytes"
	"fmt"
	"sort"
)

// Document is a mutable ordered map that compactly represents a BSON document.
type Document struct {
	// The default behavior or Append, Prepend, and Replace is to panic on the
	// insertion of a nil element. Setting IgnoreNilInsert to true will instead
	// silently ignore any nil parameters to these methods.
	IgnoreNilInsert bool
	elems           []*Element
	index           []uint32
}

// Element represents a BSON element, i.e. key-value pair of a BSON document.
type Element struct {
	value *Value
}

// Value represents a BSON value. It can be obtained as part of a bson.Element or created for use
// in a bson.Array with the bson.VC constructors.
type Value struct {
	// NOTE: For subdocuments, arrays, and code with scope, the data slice of
	// bytes may contain just the key, or the key and the code in the case of
	// code with scope. If this is the case, the start will be 0, the value will
	// be the length of the slice, and d will be non-nil.

	// start is the offset into the data slice of bytes where this element
	// begins.
	start uint32
	// offset is the offset into the data slice of bytes where this element's
	// value begins.
	offset uint32
	// data is a potentially shared slice of bytes that contains the actual
	// element. Most of the methods of this type directly index into this slice
	// of bytes.
	data []byte

	d *Document
}

type DocumentCodec struct{}

func (c *DocumentCodec) Decode(reg *CodecRegistry, vr ValueReader, v interface{}) error {
	var target *Document
	var ok bool
	if target, ok = v.(*Document); !ok {
		return fmt.Errorf("%T can only be used to decode *bson2.Document", c)
	}

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

		headerSize := 1 + len(name) + 1
		value := make([]byte, evr.Size()+headerSize)
		value[0] = byte(evr.Type())
		copy(value[1:], []byte(name))
		value[headerSize-1] = 0

		if err = evr.ReadBytes(value[headerSize:]); err != nil {
			return err
		}

		elem := &Element{&Value{
			start:  0,
			offset: uint32(headerSize),
			data:   value,
		}}

		nameBytes := []byte(name)

		target.elems = append(target.elems, elem)
		i := sort.Search(len(target.index), func(i int) bool {
			return bytes.Compare(nameBytes, elem.value.data[elem.value.start+1:elem.value.offset]) >= 0
		})
		if i < len(target.index) {
			target.index = append(target.index, 0)
			copy(target.index[i+1:], target.index[i:])
			target.index[i] = uint32(len(target.elems) - 1)
		} else {
			target.index = append(target.index, uint32(len(target.elems)-1))
		}
	}

	return nil
}

func (c *DocumentCodec) Encode(reg *CodecRegistry, vw ValueWriter, v interface{}) error {
	return fmt.Errorf("not supported")
}
