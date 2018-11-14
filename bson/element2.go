// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bson

import (
	"fmt"
)

// Elem2 represents a BSON element.
//
// NOTE: Element cannot be the value of a map nor a property of a struct without special handling.
// The default encoders and decoders will not process Element correctly. To do so would require
// information loss since an Element contains a key, but the keys used when encoding a struct are
// the struct field names. Instead of using an Element, use a Value as a value in a map or a
// property of a struct.
type Elem2 struct {
	Key   string
	Value Val2
}

func (e Elem2) String() string {
	// TODO(GODRIVER-612): When bsoncore has appenders for extended JSON use that here.
	return fmt.Sprintf(`bson.Element{"%s": %v}`, e.Key, e.Value)
}
