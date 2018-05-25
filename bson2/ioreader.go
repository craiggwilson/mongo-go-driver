package bson2

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

func NewDocumentReaderFromIO(r io.Reader) (DocumentReader, error) {
	vr, err := NewValueReaderFromIO(r, TypeDocument)
	if err != nil {
		return nil, err
	}
	return vr.ReadDocument()
}

func NewValueReaderFromIO(r io.Reader, t Type) (ValueReader, error) {
	ioReader := &ioReader{
		r:         newTrackingReader(r),
		valueType: t,
	}

	err := ioReader.setValueSize()
	if err != nil {
		return nil, err
	}

	return ioReader, nil
}

type mode byte

const (
	documentMode mode = iota
	arrayMode
)

type ioReader struct {
	r *trackingReader

	// state
	docStartPositionStack []int
	docSizeStack          []int
	currentDepth          int
	onElement             bool
	valueSize             int
	valueType             Type
}

func (r *ioReader) ensureValue(t Type) error {
	if r.valueType != t {
		return r.wrapError(newErrValueType(r.valueType, t))
	}
	if r.onElement {
		return r.wrapError(errNotValue)
	}
	return nil
}

func (r *ioReader) setValueSize() error {
	switch r.valueType {
	case TypeBoolean:
		r.valueSize = 1
	case TypeDocument:
		sizeBytes, err := r.r.peekBytes(4)
		if err != nil {
			return err
		}

		r.valueSize = int(binary.LittleEndian.Uint32(sizeBytes))
	case TypeInt32:
		r.valueSize = 4
	case TypeInt64:
		r.valueSize = 8
	case TypeString:
		sizeBytes, err := r.r.peekBytes(4)
		if err != nil {
			return err
		}

		r.valueSize = int(binary.LittleEndian.Uint32(sizeBytes)) + 4 // size is not included in string's size
	default:
		return fmt.Errorf("unsupported bson type %v", r.valueType)
	}

	return nil
}

func (r *ioReader) ReadArray() (ArrayReader, error) {
	return nil, fmt.Errorf("unsupported array")
}

func (r *ioReader) ReadBoolean() (bool, error) {
	if err := r.ensureValue(TypeBoolean); err != nil {
		return false, err
	}

	b, err := r.r.readByte()
	if err != nil {
		return false, r.wrapError(err)
	}

	r.onElement = true

	if b > 1 {
		return false, r.wrapError(errInvalidValue(fmt.Sprintf("invalid byte for boolean, %s", b)))
	}

	return b == 1, nil
}

func (r *ioReader) ReadBytes(buf []byte) error {
	if r.onElement {
		return r.wrapError(errNotValue)
	}

	if len(buf) != r.valueSize {
		return r.wrapError(fmt.Errorf("buffer must be of size %d, but got %d", r.valueSize, len(buf)))
	}

	r.onElement = true
	return r.r.readBytes(buf)
}

func (r *ioReader) ReadDocument() (DocumentReader, error) {
	if err := r.ensureValue(TypeDocument); err != nil {
		return nil, err
	}

	r.currentDepth++
	r.docStartPositionStack = append(r.docStartPositionStack, r.r.position)

	if err := r.r.fillTemp(4); err != nil {
		return nil, r.wrapError(err)
	}

	size := int(binary.LittleEndian.Uint32(r.r.temp[:4]))
	r.docSizeStack = append(r.docSizeStack, size)
	r.onElement = true
	return r, nil
}

func (r *ioReader) ReadElement() (string, ValueReader, error) {
	if !r.onElement {
		return "", nil, r.wrapError(errNotElement)
	}

	t, err := r.r.readByte()
	if err != nil {
		return "", nil, r.wrapError(err)
	}

	if t == 0 {
		r.currentDepth--
		startPosition := r.docStartPositionStack[r.currentDepth]
		size := r.docSizeStack[r.currentDepth]
		if r.r.position-startPosition != size {
			// TODO: use start position for this error report
			return "", nil, r.wrapError(errInvalidDocumentLength)
		}

		r.docStartPositionStack = r.docStartPositionStack[:r.currentDepth]
		r.docSizeStack = r.docSizeStack[:r.currentDepth]
		return "", nil, EOD
	}

	nameBytes, err := r.r.readBytesDelim(0)
	if err != nil {
		return "", nil, r.wrapError(err)
	}

	r.onElement = false
	r.valueType = Type(t)
	if err = r.setValueSize(); err != nil {
		return "", nil, err
	}
	return string(nameBytes[:len(nameBytes)-1]), r, nil
}

func (r *ioReader) ReadInt32() (int32, error) {
	if err := r.ensureValue(TypeInt32); err != nil {
		return 0, err
	}

	if err := r.r.fillTemp(4); err != nil {
		return 0, err
	}

	r.onElement = true
	return int32(binary.LittleEndian.Uint32(r.r.temp[:4])), nil
}

func (r *ioReader) ReadInt64() (int64, error) {
	if err := r.ensureValue(TypeInt64); err != nil {
		return 0, err
	}

	if err := r.r.fillTemp(8); err != nil {
		return 0, err
	}

	r.onElement = true
	return int64(binary.LittleEndian.Uint64(r.r.temp[:8])), nil
}

func (r *ioReader) ReadString() (string, error) {
	if err := r.ensureValue(TypeString); err != nil {
		return "", err
	}

	if err := r.r.fillTemp(r.valueSize); err != nil {
		return "", err
	}

	if r.r.temp[r.valueSize-1] != 0 {
		return "", r.wrapError(fmt.Errorf("string not terminated by NUL"))
	}

	r.onElement = true
	return string(r.r.temp[4 : r.valueSize-1]), nil
}

func (r *ioReader) Size() int {
	return r.valueSize
}

func (r *ioReader) Skip() error {
	if r.onElement {
		return r.wrapError(errNotValue)
	}

	return r.r.skip(r.valueSize)
}

func (r *ioReader) Type() Type {
	return r.valueType
}

func (r *ioReader) wrapError(err error) error {
	return fmt.Errorf("position %d: %s", r.r.position, err)
}

func newTrackingReader(r io.Reader) *trackingReader {
	return &trackingReader{
		r: bufio.NewReaderSize(r, 512),
	}
}

type trackingReader struct {
	r *bufio.Reader

	position int

	temp []byte
}

func (r *trackingReader) fillTemp(size int) error {
	if len(r.temp) < size {
		r.temp = make([]byte, size)
	}

	return r.readBytes(r.temp[:size])
}

func (r *trackingReader) peekBytes(size int) ([]byte, error) {
	return r.r.Peek(4)
}

func (r *trackingReader) readByte() (byte, error) {
	b, err := r.r.ReadByte()
	if err != nil {
		return 0, err
	}

	r.position++
	return b, nil
}

func (r *trackingReader) readBytes(buf []byte) error {
	n, err := io.ReadFull(r.r, buf)
	r.position += n
	return err
}

func (r *trackingReader) readBytesDelim(delim byte) ([]byte, error) {
	result, err := r.r.ReadBytes(delim)
	r.position += len(result)
	return result, err
}

func (r *trackingReader) skip(size int) error {
	n, err := r.r.Discard(size)
	r.position += n
	return err
}
