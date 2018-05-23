package bson2

import (
	"bytes"
	"testing"
)

func TestReader(t *testing.T) {
	ioReader := bytes.NewBuffer([]byte{0x05, 0x00, 0x00, 0x00, 0x00})

	bsonReader := NewReaderFromIO(ioReader)

	name, err := bsonReader.ReadName()
}
