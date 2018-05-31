package bson2

import (
	"encoding/binary"
	"fmt"
	"io"
)

type ArrayReader ValueReader

type DocumentReader interface {
	ReadElement() (string, ValueReader, error)
}

type ValueReader interface {
	// ReadBytes gets the bytes representing the value.
	ReadBytes([]byte) error
	ReadSlice() ([]byte, error)
	Size() int
	Type() Type

	ReadArray() (ArrayReader, error)
	ReadBoolean() (bool, error)
	ReadDocument() (DocumentReader, error)
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadString() (string, error)
	Skip() error
}

func NewDocumentReader(input []byte) (DocumentReader, error) {
	vr, err := NewValueReader(input, TypeDocument)
	if err != nil {
		return nil, err
	}
	return vr.ReadDocument()
}

func NewValueReader(input []byte, t Type) (ValueReader, error) {
	r := &reader{
		data: input,
	}

	r.state.valueType = t
	err := r.setValueSize()
	if err != nil {
		return nil, err
	}

	return r, nil
}

type mode byte

const (
	documentMode mode = iota
	arrayMode
)

type reader struct {
	data []byte
	pos  int

	state readerState
}

type readerState struct {
	prev *readerState

	docStartPos int
	docSize     int
	depth       int
	onElement   bool
	valueSize   int
	valueType   Type
}

func (s *readerState) push() {
	if s.depth > 0 {
		s.prev = &readerState{
			prev:        s.prev,
			docStartPos: s.docStartPos,
			docSize:     s.docSize,
			depth:       s.depth,
			onElement:   s.onElement,
			valueSize:   s.valueSize,
			valueType:   s.valueType,
		}
	}
	s.depth++
}

func (s *readerState) pop() {
	if s.prev != nil {
		s.docStartPos = s.prev.docStartPos
		s.docSize = s.prev.docSize
		s.depth = s.prev.depth
		s.onElement = s.prev.onElement
		s.valueSize = s.prev.valueSize
		s.valueType = s.prev.valueType
		s.prev = s.prev.prev
	} else {
		s.depth--
	}
}

func (r *reader) ensureValue(t Type) error {
	if r.state.valueType != t {
		return r.wrapError(newErrValueType(r.state.valueType, t))
	}
	if r.state.onElement {
		return r.wrapError(errNotValue)
	}
	return nil
}

func (r *reader) setValueSize() error {
	switch r.state.valueType {
	case TypeBoolean:
		r.state.valueSize = 1
	case TypeDocument:
		sizeBytes, err := r.peekBytes(4)
		if err != nil {
			return err
		}

		r.state.valueSize = int(binary.LittleEndian.Uint32(sizeBytes))
	case TypeInt32:
		r.state.valueSize = 4
	case TypeInt64:
		r.state.valueSize = 8
	case TypeString:
		sizeBytes, err := r.peekBytes(4)
		if err != nil {
			return err
		}

		r.state.valueSize = int(binary.LittleEndian.Uint32(sizeBytes)) + 4 // size is not included in string's size
	default:
		return fmt.Errorf("unsupported bson type %v", r.state.valueType)
	}

	return nil
}

func (r *reader) ReadArray() (ArrayReader, error) {
	return nil, fmt.Errorf("unsupported array")
}

func (r *reader) ReadBoolean() (bool, error) {
	if err := r.ensureValue(TypeBoolean); err != nil {
		return false, err
	}

	b, err := r.readByte()
	if err != nil {
		return false, r.wrapError(err)
	}

	r.state.onElement = true

	if b > 1 {
		return false, r.wrapError(errInvalidValue(fmt.Sprintf("invalid byte for boolean, %s", b)))
	}

	return b == 1, nil
}

func (r *reader) ReadBytes(buf []byte) error {
	if r.state.onElement {
		return r.wrapError(errNotValue)
	}

	if len(buf) != r.state.valueSize {
		return r.wrapError(fmt.Errorf("buffer must be of size %d, but got %d", r.state.valueSize, len(buf)))
	}

	r.state.onElement = true
	data, err := r.readBytes(len(buf))
	if err != nil {
		return r.wrapError(err)
	}
	copy(buf, data)
	return nil
}

func (r *reader) ReadSlice() ([]byte, error) {
	if r.state.onElement {
		return nil, r.wrapError(errNotValue)
	}

	r.state.onElement = true
	data, err := r.readBytes(r.state.valueSize)
	if err != nil {
		return nil, r.wrapError(err)
	}
	return data, nil
}

func (r *reader) ReadDocument() (DocumentReader, error) {
	if err := r.ensureValue(TypeDocument); err != nil {
		return nil, err
	}

	// save current state if necessary
	r.state.push()
	r.state.docStartPos = r.pos

	data, err := r.readBytes(4)
	if err != nil {
		return nil, r.wrapError(err)
	}

	r.state.docSize = int(binary.LittleEndian.Uint32(data))
	r.state.onElement = true
	return r, nil
}

func (r *reader) ReadElement() (string, ValueReader, error) {
	if !r.state.onElement {
		return "", nil, r.wrapError(errNotElement)
	}

	t, err := r.readByte()
	if err != nil {
		return "", nil, r.wrapError(err)
	}

	if t == 0 {
		if r.pos-r.state.docStartPos != r.state.docSize {
			// TODO: use start position for this error report
			return "", nil, r.wrapError(errInvalidDocumentLength)
		}

		r.state.pop()
		r.state.onElement = true
		return "", nil, EOD
	}

	nameBytes, err := r.readBytesDelim(0)
	if err != nil {
		return "", nil, r.wrapError(err)
	}

	r.state.onElement = false
	r.state.valueType = Type(t)
	if err = r.setValueSize(); err != nil {
		return "", nil, err
	}
	return string(nameBytes[:len(nameBytes)-1]), r, nil
}

func (r *reader) ReadInt32() (int32, error) {
	if err := r.ensureValue(TypeInt32); err != nil {
		return 0, r.wrapError(err)
	}

	data, err := r.readBytes(4)
	if err != nil {
		return 0, r.wrapError(err)
	}

	r.state.onElement = true
	return int32(binary.LittleEndian.Uint32(data)), nil
}

func (r *reader) ReadInt64() (int64, error) {
	if err := r.ensureValue(TypeInt64); err != nil {
		return 0, r.wrapError(err)
	}

	data, err := r.readBytes(8)
	if err != nil {
		return 0, r.wrapError(err)
	}

	r.state.onElement = true
	return int64(binary.LittleEndian.Uint64(data)), nil
}

func (r *reader) ReadString() (string, error) {
	if err := r.ensureValue(TypeString); err != nil {
		return "", r.wrapError(err)
	}

	data, err := r.readBytes(r.state.valueSize)
	if err != nil {
		return "", r.wrapError(err)
	}

	if data[r.state.valueSize-1] != 0 {
		return "", r.wrapError(fmt.Errorf("string not terminated by NUL"))
	}

	r.state.onElement = true
	return string(data[4 : len(data)-1]), nil
}

func (r *reader) Size() int {
	return r.state.valueSize
}

func (r *reader) Skip() error {
	if r.state.onElement {
		return r.wrapError(errNotValue)
	}

	return r.skip(r.state.valueSize)
}

func (r *reader) Type() Type {
	return r.state.valueType
}

func (r *reader) wrapError(err error) error {
	return fmt.Errorf("position %d: %s", r.pos, err)
}

func (r *reader) String() string {
	return fmt.Sprintf("pos=%d, remaining=%v", r.pos, r.data[r.pos:])
}

func (r *reader) peekBytes(size int) ([]byte, error) {
	if r.pos+size > len(r.data) {
		return nil, io.EOF
	}

	return r.data[r.pos : r.pos+size], nil
}

func (r *reader) readByte() (byte, error) {
	if r.pos+1 > len(r.data) {
		return 0, io.EOF
	}

	r.pos++
	return r.data[r.pos-1], nil
}

func (r *reader) readBytes(size int) ([]byte, error) {
	if r.pos+size > len(r.data) {
		return nil, io.EOF
	}
	r.pos += size
	return r.data[r.pos-size : r.pos], nil
}

func (r *reader) readBytesDelim(delim byte) ([]byte, error) {
	start := r.pos
	for ; r.pos < len(r.data) && r.data[r.pos] != delim; r.pos++ {
	}
	if r.pos >= len(r.data) {
		r.pos = start
		return nil, io.EOF
	} else if r.data[r.pos] == delim {
		r.pos++
	}
	return r.data[start:r.pos], nil
}

func (r *reader) skip(size int) error {
	if r.pos+size > len(r.data) {
		return io.EOF
	}
	r.pos += size
	return nil
}
