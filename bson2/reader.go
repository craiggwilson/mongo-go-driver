package bson2

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

func NewDocumentReaderFromIO(r io.Reader) (DocumentReader, error) {
	vr := NewValueReader(r, TypeDocument)
	return vr.ReadDocument()
}

func NewValueReader(r io.Reader, t Type) ValueReader {
	return &ioReader{
		r:         newTrackingReader(r),
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

	// ReadBytes gets the bytes representing the value.
	ReadBytes() ([]byte, error)

	ReadBoolean() (bool, error)
	ReadDocument() (DocumentReader, error)
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadString() (string, error)
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
	currentElementName    string
	onElement             bool
	valueType             Type
}

func (r *ioReader) readValueBytes(t Type) ([]byte, error) {
	if r.valueType != t {
		return nil, r.wrapError(newErrValueType(r.valueType, t))
	}

	return r.ReadBytes()
}

func (r *ioReader) ReadArray() (ArrayReader, error) {
	if r.valueType != TypeArray {
		return nil, r.wrapError(newErrValueType(r.valueType, TypeArray))
	}
	if r.onElement {
		return nil, r.wrapError(errNotValue)
	}
	r.onElement = true

	r.currentDepth++

	// TODO: do something with the size
	_, err := r.r.readInt()
	if err != nil {
		return nil, r.wrapError(err)
	}

	return r, nil
}

func (r *ioReader) ReadBoolean() (bool, error) {
	bytes, err := r.readValueBytes(TypeBoolean)
	if err != nil {
		return false, r.wrapError(err)
	}

	if bytes[0] == 0 {
		return false, nil
	} else if bytes[0] == 1 {
		return true, nil
	}

	return false, r.wrapError(errInvalidValue(fmt.Sprintf("invalid byte for boolean, %s", bytes[0])))
}

func (r *ioReader) ReadBytes() ([]byte, error) {
	if r.onElement {
		return nil, r.wrapError(errNotValue)
	}
	r.onElement = true

	switch r.valueType {
	case TypeBoolean:
		b, err := r.r.readByte()
		if err != nil {
			return nil, r.wrapError(err)
		}
		return []byte{b}, nil
	case TypeDocument:
		sizeBytes, err := r.r.readBytes(4)
		if err != nil {
			return nil, r.wrapError(err)
		}

		size := int(binary.LittleEndian.Uint32(sizeBytes))
		result := append([]byte{}, sizeBytes...)

		bytes, err := r.r.readBytes(size - 4)
		if err != nil {
			return nil, r.wrapError(err)
		}

		result = append(result, bytes...)

		return result, nil
	case TypeInt32:
		result, err := r.r.readBytes(4)
		if err != nil {
			return nil, r.wrapError(err)
		}

		return result, nil
	case TypeInt64:
		result, err := r.r.readBytes(8)
		if err != nil {
			return nil, r.wrapError(err)
		}

		return result, nil
	case TypeString:
		sizeBytes, err := r.r.readBytes(4)
		if err != nil {
			return nil, r.wrapError(err)
		}

		size := int(binary.LittleEndian.Uint32(sizeBytes))

		result := append([]byte{}, sizeBytes...)

		bytes, err := r.r.readBytes(size)
		if err != nil {
			return nil, r.wrapError(err)
		}

		result = append(result, bytes...)

		return result, nil
	default:
		return nil, fmt.Errorf("unsupported bson type %v", r.valueType)
	}
}

func (r *ioReader) ReadDocument() (DocumentReader, error) {
	if r.valueType != TypeDocument {
		return nil, r.wrapError(newErrValueType(r.valueType, TypeDocument))
	}
	if r.onElement {
		return nil, r.wrapError(errNotValue)
	}
	r.onElement = true

	r.currentDepth++

	r.docStartPositionStack = append(r.docStartPositionStack, r.r.position)

	size, err := r.r.readInt()
	if err != nil {
		return nil, r.wrapError(err)
	}

	r.docSizeStack = append(r.docSizeStack, size)

	return r, nil
}

func (r *ioReader) ReadElement() (string, ValueReader, error) {
	if !r.onElement {
		return "", nil, r.wrapError(errNotElement)
	}
	r.onElement = false

	t, err := r.r.readByte()
	if err != nil {
		return "", nil, r.wrapError(err)
	}

	if t == 0 {
		r.currentDepth--
		startPosition := r.docStartPositionStack[r.currentDepth]
		size := r.docSizeStack[r.currentDepth]
		if r.r.position-startPosition != size {
			return "", nil, r.wrapError(errInvalidDocumentLength)
		}

		r.docStartPositionStack = r.docStartPositionStack[:r.currentDepth]
		r.docSizeStack = r.docSizeStack[:r.currentDepth]
		r.onElement = true // go back to outer document
		return "", nil, EOD
	}

	nameBytes, err := r.r.readBytesDelim(0)
	if err != nil {
		return "", nil, r.wrapError(err)
	}

	r.valueType = Type(t)

	r.currentElementName = string(nameBytes[:len(nameBytes)-1])

	return r.currentElementName, r, nil
}

func (r *ioReader) ReadInt32() (int32, error) {
	bytes, err := r.readValueBytes(TypeInt32)
	if err != nil {
		return 0, r.wrapError(err)
	}

	return int32(binary.LittleEndian.Uint32(bytes)), nil
}

func (r *ioReader) ReadInt64() (int64, error) {
	bytes, err := r.readValueBytes(TypeInt64)
	if err != nil {
		return 0, r.wrapError(err)
	}

	return int64(binary.LittleEndian.Uint64(bytes)), nil
}

func (r *ioReader) ReadString() (string, error) {
	bytes, err := r.readValueBytes(TypeString)
	if err != nil {
		return "", r.wrapError(err)
	}

	return string(bytes), nil
}

func (r *ioReader) Type() Type {
	return r.valueType
}

func (r *ioReader) wrapError(err error) error {
	return fmt.Errorf("position %d: %s", r.r.position, err)
}

func newTrackingReader(r io.Reader) *trackingReader {
	return &trackingReader{
		r: bufio.NewReader(r),
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

	n, err := io.ReadFull(r.r, r.temp[:size])
	r.position += n

	return err
}

func (r *trackingReader) readByte() (byte, error) {
	b, err := r.r.ReadByte()
	if err != nil {
		return 0, err
	}

	r.position++

	return b, nil
}

func (r *trackingReader) readBytes(size int) ([]byte, error) {
	err := r.fillTemp(size)
	if err != nil {
		return nil, err
	}

	return r.temp[:size], nil
}

func (r *trackingReader) readBytesDelim(delim byte) ([]byte, error) {
	result, err := r.r.ReadBytes(delim)
	if err != nil {
		return nil, err
	}

	r.position += len(result)

	return result, nil
}

func (r *trackingReader) readInt() (int, error) {
	bytes, err := r.readBytes(4)
	if err != nil {
		return 0, err
	}

	return int(binary.LittleEndian.Uint32(bytes)), nil
}
