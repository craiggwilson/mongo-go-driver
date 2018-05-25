package bson2

type ArrayReader ValueReader

type DocumentReader interface {
	ReadElement() (string, ValueReader, error)
}

type ValueReader interface {
	// ReadBytes gets the bytes representing the value.
	ReadBytes([]byte) error
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
