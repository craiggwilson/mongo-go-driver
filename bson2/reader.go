package bson2

import (
	"encoding/binary"
	"io"
)

func NewReaderFromIO(r io.Reader) DocumentReader {
	return &ioReader{
		r:         r,
		onElement: false,
		valueType: TypeDocument,
	}
}

type ArrayReader ValueReader

type DocumentReader interface {
	ReadElement() (string, ValueReader, error)
}

type ValueReader interface {
	Type() Type

	ReadArray() (ArrayReader, error)
	ReadDocument() (DocumentReader, error)
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadString() (string, error)
}

type ioReader struct {
	r io.Reader

	limits []int

	onElement bool
	valueType Type

	temp []byte
}

func (r *ioReader) Type() Type {
	return r.valueType
}

func (r *ioReader) ReadElement() (string, ValueReader, error) {
	if !r.onElement {
		return "", nil, ErrNotElement
	}

	err := r.fillTemp(1)
	if err != nil {
		return "", nil, err
	}

	t := r.temp[0]

	name, err := r.buf.ReadString(0)
	if err != nil {
		return "", nil, err
	}

	r.valueType = Type(t)

	r.onElement = false

	return name, r, nil
}

func (r *ioReader) ReadArray() (ArrayReader, error) {
	if r.onElement {
		return nil, ErrNotValue
	}
	if r.valueType != TypeArray {
		return nil, NewErrValueType(r.valueType, TypeArray)
	}

	return r, nil
}

func (r *ioReader) ReadDocument() (DocumentReader, error) {
	if r.onElement {
		return nil, ErrNotValue
	}
	if r.valueType != TypeDocument {
		return nil, NewErrValueType(r.valueType, TypeDocument)
	}

	return r, nil
}

func (r *ioReader) ReadInt32() (int32, error) {
	if r.onElement {
		return 0, ErrNotValue
	}
	if r.valueType != TypeInt32 {
		return 0, NewErrValueType(r.valueType, TypeInt32)
	}

	err := r.fillTemp(4)
	if err != nil {
		return 0, err
	}

	r.onElement = true

	return int32(binary.LittleEndian.Uint32(r.temp)), nil
}

func (r *ioReader) ReadInt64() (int64, error) {
	if r.onElement {
		return 0, ErrNotValue
	}
	if r.valueType != TypeInt64 {
		return 0, NewErrValueType(r.valueType, TypeInt64)
	}

	err := r.fillTemp(8)
	if err != nil {
		return 0, err
	}

	r.onElement = true

	return int64(binary.LittleEndian.Uint64(r.temp)), nil
}

func (r *ioReader) ReadString() (string, error) {
	if r.onElement {
		return "", ErrNotValue
	}
	if r.valueType != TypeString {
		return "", NewErrValueType(r.valueType, TypeString)
	}

	err := r.fillTemp(4)
	if err != nil {
		return "", err
	}

	size := int(binary.LittleEndian.Uint32(r.temp))

	err = r.fillTemp(size)
	if err != nil {
		return "", err
	}

	b, err := r.buf.ReadByte()
	if err != nil {
		return "", err
	}

	if b != 0 {
		return "", ErrInvalidValue("String value was not terminated by a NUL byte")
	}

	r.onElement = true

	return string(r.temp[:size]), nil
}

func (r *ioReader) fillTemp(size int) error {
	if len(r.temp) < size {
		r.temp = make([]byte, size, size)
	}

	_, err := io.ReadFull(r.buf, r.temp[:size])

	return err
}

func (r *ioReader) readCString() (string, error) {

	var temp []byte

	for {
		r.r.Read
	}

	return nil
}
