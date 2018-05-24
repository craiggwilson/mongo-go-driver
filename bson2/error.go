package bson2

import (
	"errors"
	"fmt"
)

var EOA = errors.New("end of array")
var EOD = errors.New("end of document")

var errInvalidDocumentLength = errors.New("invalid document length")
var errNotElement = errors.New("not positioned on an element")
var errNotValue = errors.New("not positioned on a value")

func newErrValueType(currentType, attemptedType Type) error {
	return errValueType{
		currentType:   currentType,
		attemptedType: attemptedType,
	}
}

type errValueType struct {
	currentType   Type
	attemptedType Type
}

func (e errValueType) Error() string {
	return fmt.Sprintf("positioned on %s, but attempted to read %s", e.currentType, e.attemptedType)
}

type errInvalidValue string

func (e errInvalidValue) Error() string {
	return string(e)
}
