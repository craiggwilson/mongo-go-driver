// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bson

import (
	"github.com/mongodb/mongo-go-driver/bson/bsontype"
)

type Val2 interface {
	Type() bsontype.Type

	MarshalBSONValue() (bsontype.Type, []byte, error)
}

type Int32Value int32

func (Int32Value) Type() bsontype.Type {
	return bsontype.Int32
}

func (i32 Int32Value) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsontype.Int32, []byte{0,0,0,0}, nil
}