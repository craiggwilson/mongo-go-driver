package bson2

import (
	"errors"
	"fmt"
)

var ErrNotElement = errors.New("not positioned on an element")
var ErrNotValue = errors.New("not positioned on a value")

func NewErrValueType(currentType, attemptedType Type) error {
	return ErrValueType{
		currentType:   currentType,
		attemptedType: attemptedType,
	}
}

type ErrValueType struct {
	currentType   Type
	attemptedType Type
}

func (e ErrValueType) Error() string {
	return fmt.Sprintf("positioned on %s, but attempted to read %s", e.currentType, e.attemptedType)
}

type ErrInvalidValue string

func (e ErrInvalidValue) Error() string {
	return string(e)
}
