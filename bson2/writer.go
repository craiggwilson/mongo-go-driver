package bson2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type DocumentWriter interface {
	WriteElement(string) (ValueWriter, error)
	WriteEndDocument() error
}

type ValueWriter interface {
	WriteTo(w io.Writer) (int64, error)

	WriteBoolean(bool) error
	WriteDocument() (DocumentWriter, error)
	WriteInt32(int32) error
	WriteInt64(int64) error
	WriteString(string) error
}

func NewValueWriter() ValueWriter {
	w := &writer{}
	w.state.data = &bytes.Buffer{}
	return w
}

type writer struct {
	pos int

	state writerState
}

type writerState struct {
	prev *writerState

	elementName    string
	hasElementName bool
	docStartPos    int
	depth          int
	onElement      bool
	data           *bytes.Buffer
}

func (s *writerState) push() {
	s.prev = &writerState{
		prev:        s.prev,
		depth:       s.depth,
		docStartPos: s.docStartPos,
		onElement:   s.onElement,
		data:        s.data,
	}
	s.data = &bytes.Buffer{}
	s.depth++
}

func (s *writerState) pop() {
	s.docStartPos = s.prev.docStartPos
	s.depth = s.prev.depth
	s.onElement = s.prev.onElement
	s.data = s.prev.data
	s.prev = s.prev.prev
}

func (w *writer) writeValueType(t Type) error {
	if w.state.onElement {
		return errNotValue
	}

	if w.state.depth > 0 {
		if err := w.writeByte(byte(t)); err != nil {
			return err
		}
	}

	if w.state.hasElementName {
		if err := w.writeBytes([]byte(w.state.elementName)); err != nil {
			return err
		}
		if err := w.writeByte(0); err != nil {
			return err
		}
	}

	return nil
}

func (w *writer) WriteBoolean(v bool) error {
	if err := w.writeValueType(TypeBoolean); err != nil {
		return w.wrapError(err)
	}

	if v {
		if err := w.writeByte(1); err != nil {
			return w.wrapError(err)
		}
	} else {
		if err := w.writeByte(0); err != nil {
			return w.wrapError(err)
		}
	}
	w.state.onElement = true
	return nil
}

func (w *writer) WriteBytes(v []byte) error {
	// if err := w.ensureValue(); err != nil {
	// 	return w.wrapError(err)
	// }

	if err := w.writeBytes(v); err != nil {
		return w.wrapError(err)
	}

	w.state.onElement = true
	return nil
}

func (w *writer) WriteDocument() (DocumentWriter, error) {
	if err := w.writeValueType(TypeDocument); err != nil {
		return nil, w.wrapError(err)
	}

	w.state.push()
	w.state.onElement = true
	w.pos += 4 // accounting for unwritten size bytes

	return w, nil
}

func (w *writer) WriteElement(name string) (ValueWriter, error) {
	if !w.state.onElement {
		return nil, w.wrapError(errNotElement)
	}

	w.state.elementName = name
	w.state.hasElementName = true

	w.state.onElement = false
	return w, nil
}

func (w *writer) WriteEndDocument() error {
	if err := w.writeByte(0); err != nil {
		return err
	}

	docData := w.state.data
	w.state.pop()

	docSize := docData.Len() + 4
	docSizeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(docSizeBytes, uint32(docSize))
	if err := w.writeBytes(docSizeBytes); err != nil {
		return w.wrapError(err)
	}
	if _, err := docData.WriteTo(w.state.data); err != nil {
		return w.wrapError(err)
	}

	return nil
}

func (w *writer) WriteInt32(v int32) error {
	if err := w.writeValueType(TypeInt32); err != nil {
		return w.wrapError(err)
	}

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(v))
	if err := w.writeBytes(b); err != nil {
		return w.wrapError(err)
	}

	w.state.onElement = true
	return nil
}

func (w *writer) WriteInt64(v int64) error {
	if err := w.writeValueType(TypeInt64); err != nil {
		return w.wrapError(err)
	}

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(v))
	if err := w.writeBytes(b); err != nil {
		return w.wrapError(err)
	}

	w.state.onElement = true
	return nil
}

func (w *writer) WriteString(v string) error {
	if err := w.writeValueType(TypeString); err != nil {
		return w.wrapError(err)
	}

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(len(v)+1))
	if err := w.writeBytes(b); err != nil {
		return w.wrapError(err)
	}

	if err := w.writeBytes([]byte(v)); err != nil {
		return w.wrapError(err)
	}

	if err := w.writeByte(0); err != nil {
		return w.wrapError(err)
	}

	w.state.onElement = true
	return nil
}

func (w *writer) WriteTo(out io.Writer) (int64, error) {
	return w.state.data.WriteTo(out)
}

func (w *writer) wrapError(err error) error {
	return fmt.Errorf("position %d: %s", w.pos, err)
}

func (w *writer) writeByte(b byte) error {
	if err := w.state.data.WriteByte(b); err != nil {
		return err
	}

	w.pos++
	return nil
}

func (w *writer) writeBytes(b []byte) error {
	n, err := w.state.data.Write(b)
	w.pos += n
	return err
}
