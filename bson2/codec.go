package bson2

type Decoder interface {
	Decode(r ValueReader, value interface{}) (interface{}, error)
}
