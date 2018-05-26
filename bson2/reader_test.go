package bson2

import (
	"encoding/hex"
	"errors"
	"testing"
)

func TestReader(t *testing.T) {
	tests := []struct {
		name  string
		input string
		err   error
	}{
		{
			name:  "empty document",
			input: "0500000000",
		},
		{
			name:  "empty sub document",
			input: "0D000000037800050000000000",
		},
		{
			name:  "empty string key subdoc",
			input: "150000000378000D00000002000200000062000000",
		},
		{
			name:  "single character key subdoc",
			input: "160000000378000E0000000261000200000062000000",
		},
		{
			name:  "top document length too long: eats terminator",
			input: "0600000000",
			err:   errors.New("position 5: invalid document length"),
		},
		{
			name:  "sub document length too long: eats outer terminator",
			input: "1800000003666F6F000F0000001062617200FFFFFF7F0000",
			err:   errors.New("position 23: invalid document length"),
		},
		{
			name:  "sub document length too short: leaks terminator",
			input: "1500000003666F6F000A0000000862617200010000",
			err:   errors.New("position 20: invalid document length"),
		},
		{
			name:  "sub document invalid: bad string length in field",
			input: "1C00000003666F6F001200000002626172000500000062617A000000",
			err:   errors.New("position 28: invalid document length"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input, err := hex.DecodeString(test.input)
			if err != nil {
				t.Fatal(err)
			}

			bsonReader, err := NewDocumentReader(input)
			if err != nil {
				t.Fatal(err)
			}

			err = validateBsonDocument(t, bsonReader)
			if err == nil && test.err != nil {
				t.Fatalf("expected error '%s', but got none", test.err)
			} else if err != nil && test.err == nil {
				t.Fatalf("unexpected error '%s'", err)
			} else if err != nil && test.err != nil && err.Error() != test.err.Error() {
				t.Fatalf("expected error '%s', but got '%s'", test.err, err)
			}
		})
	}
}

func validateBsonDocument(t *testing.T, bsonReader DocumentReader) error {
	for {
		_, vr, err := bsonReader.ReadElement()
		if err == EOD { // end of document, expected
			break
		}

		if err != nil {
			return err
		}

		switch vr.Type() {
		case TypeDocument:
			dr, err := vr.ReadDocument()
			if err != nil {
				return err
			}
			err = validateBsonDocument(t, dr)
			if err != nil {
				return err
			}
		default:
			bytes := make([]byte, vr.Size())
			err := vr.ReadBytes(bytes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
