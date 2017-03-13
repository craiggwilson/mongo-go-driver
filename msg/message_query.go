package msg

import "fmt"

// Query is a message sent to the server.
type Query struct {
	ReqID                int32
	Flags                QueryFlags
	FullCollectionName   string
	NumberToSkip         int32
	NumberToReturn       int32
	Query                interface{}
	ReturnFieldsSelector interface{}
}

// RequestID gets the request id of the message.
func (m *Query) RequestID() int32 { return m.ReqID }

// QueryFlags are the flags in a Query.
type QueryFlags int32

// QueryFlags constants.
const (
	_ QueryFlags = 1 << iota
	TailableCursor
	SlaveOK
	OplogReplay
	NoCursorTimeout
	AwaitData
	Exhaust
	Partial
)

// WrapWithMeta wraps the query with meta data.
func WrapWithMeta(r Request, meta map[string]interface{}) {
	if len(meta) > 0 {
		switch typedR := r.(type) {
		case *Query:
			typedR.Query = struct {
				Q interface{}            `bson:"$query"`
				M map[string]interface{} `bson:",inline"`
			}{
				Q: typedR.Query,
				M: meta,
			}
		default:
			panic(fmt.Sprintf("cannot wrap request(%T) with meta", r))
		}
	}
}
