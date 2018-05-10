package bson

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode2(t *testing.T) {
	registry := NewCodecRegistry()
	codec := NewCodecRegistryCodec(registry)

	t.Run("Decode to []byte", func(t *testing.T) {
		input := []byte{0x5, 0x0, 0x0, 0x0, 0x0}
		result := make([]byte, 5)
		actual, err := codec.Decode(bytes.NewBuffer(input), result)
		require.NoError(t, err)

		require.Equal(t, []byte{0x5, 0x0, 0x0, 0x0, 0x0}, result)
		require.Equal(t, result, actual)
	})
	t.Run("Decode to map", func(t *testing.T) {
		tests := []struct {
			name   string
			input  *Document
			result interface{}
			output map[string]interface{}
		}{{
			name:   "empty document",
			input:  NewDocument(),
			result: map[string]interface{}{},
			output: map[string]interface{}{},
		}, {
			name: "non-empty document into empty document",
			input: NewDocument(
				EC.String("foo", "bar"),
				EC.Int32("baz", 32),
			),
			result: map[string]interface{}{},
			output: map[string]interface{}{
				"foo": "bar",
				"baz": int32(32),
			},
		}, {
			name: "non-empty document into non-empty document",
			input: NewDocument(
				EC.String("foo", "bar"),
				EC.Int32("baz", 32),
			),
			result: map[string]interface{}{
				"foo": "overwrite",
				"bar": int32(12),
			},
			output: map[string]interface{}{
				"foo": "bar",
				"bar": int32(12),
				"baz": int32(32),
			},
		}, {
			name: "non-empty document with a subdocument",
			input: NewDocument(
				EC.String("foo", "bar"),
				EC.Int32("baz", 32),
				EC.SubDocumentFromElements(
					"sub",
					EC.Int32("foo", 13),
				),
			),
			result: map[string]interface{}{},
			output: map[string]interface{}{
				"foo": "bar",
				"baz": int32(32),
				"sub": map[string]interface{}{
					"foo": int32(13),
				},
			},
		}, {
			name: "non-empty document with a non-empty document/subdocument",
			input: NewDocument(
				EC.String("foo", "bar"),
				EC.Int32("baz", 32),
				EC.SubDocumentFromElements(
					"sub",
					EC.Int32("foo", 13),
				),
			),
			result: map[string]interface{}{
				"foo": "overwrite",
				"bar": int32(12),
				"sub": map[string]interface{}{
					"foo": int32(14),
					"bar": int32(16),
				},
			},
			output: map[string]interface{}{
				"foo": "bar",
				"bar": int32(12),
				"baz": int32(32),
				"sub": map[string]interface{}{
					"foo": int32(13),
					"bar": int32(16),
				},
			},
		}, {
			name: "non-empty document with a non-empty document/typed subdocument",
			input: NewDocument(
				EC.String("foo", "bar"),
				EC.Int32("baz", 32),
				EC.SubDocumentFromElements(
					"sub",
					EC.Int32("foo", 13),
				),
			),
			result: map[string]interface{}{
				"foo": "overwrite",
				"bar": int32(12),
				"sub": map[string]int32{
					"foo": int32(14),
					"bar": int32(16),
				},
			},
			output: map[string]interface{}{
				"foo": "bar",
				"bar": int32(12),
				"baz": int32(32),
				"sub": map[string]int32{
					"foo": int32(13),
					"bar": int32(16),
				},
			},
		}}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				var buffer bytes.Buffer
				test.input.WriteTo(&buffer)
				actual, err := codec.Decode(&buffer, test.result)
				require.NoError(t, err)

				require.Equal(t, test.output, test.result)
				require.Equal(t, test.result, actual)
			})
		}
	})
	t.Run("Decode to struct", func(t *testing.T) {
		type sub struct {
			Foo int32
			Bar int32
			Baz int32
		}
		type target struct {
			Foo string
			Bar int32
			Baz int32
			Sub sub
		}

		tests := []struct {
			name   string
			input  *Document
			result *target
			output *target
		}{{
			name:   "empty document",
			input:  NewDocument(),
			result: &target{},
			output: &target{},
		}, {
			name: "non-empty document into empty struct",
			input: NewDocument(
				EC.String("foo", "bar"),
				EC.Int32("baz", 32),
			),
			result: &target{},
			output: &target{
				Foo: "bar",
				Baz: 32,
			},
		}, {
			name: "non-empty document into non-empty struct",
			input: NewDocument(
				EC.String("foo", "bar"),
				EC.Int32("baz", 32),
			),
			result: &target{
				Foo: "overwrite",
				Bar: 12,
			},
			output: &target{
				Foo: "bar",
				Bar: 12,
				Baz: 32,
			},
		}, {
			name: "non-empty document into empty struct with empty subdocument",
			input: NewDocument(
				EC.String("foo", "bar"),
				EC.Int32("baz", 32),
				EC.SubDocumentFromElements(
					"sub",
					EC.Int32("foo", 13),
				),
			),
			result: &target{},
			output: &target{
				Foo: "bar",
				Baz: 32,
				Sub: sub{Foo: 13},
			},
		}, {
			name: "non-empty document into empty struct with non-empty subdocument",
			input: NewDocument(
				EC.String("foo", "bar"),
				EC.Int32("baz", 32),
				EC.SubDocumentFromElements(
					"sub",
					EC.Int32("foo", 13),
				),
			),
			result: &target{
				Foo: "overwrite",
				Bar: 13,
				Sub: sub{Foo: 12, Bar: 10},
			},
			output: &target{
				Foo: "bar",
				Bar: 13,
				Baz: 32,
				Sub: sub{Foo: 13, Bar: 10},
			},
		}}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				var buffer bytes.Buffer
				test.input.WriteTo(&buffer)
				actual, err := codec.Decode(&buffer, test.result)
				require.NoError(t, err)

				require.Equal(t, test.output, test.result)
				require.Equal(t, test.result, actual)
			})
		}
	})
	t.Run("Decode to *struct", func(t *testing.T) {
		type sub struct {
			Foo int32
			Bar int32
			Baz int32
		}
		type target struct {
			Foo string
			Bar int32
			Baz int32
			Sub *sub
		}

		tests := []struct {
			name   string
			input  *Document
			result *target
			output *target
		}{{
			name:   "empty document",
			input:  NewDocument(),
			result: &target{},
			output: &target{},
		}, {
			name: "non-empty document into empty struct",
			input: NewDocument(
				EC.String("foo", "bar"),
				EC.Int32("baz", 32),
			),
			result: &target{},
			output: &target{
				Foo: "bar",
				Baz: 32,
			},
		}, {
			name: "non-empty document into non-empty struct",
			input: NewDocument(
				EC.String("foo", "bar"),
				EC.Int32("baz", 32),
			),
			result: &target{
				Foo: "overwrite",
				Bar: 12,
			},
			output: &target{
				Foo: "bar",
				Bar: 12,
				Baz: 32,
			},
		}, {
			name: "non-empty document into empty struct with subdocument",
			input: NewDocument(
				EC.String("foo", "bar"),
				EC.Int32("baz", 32),
				EC.SubDocumentFromElements(
					"sub",
					EC.Int32("foo", 13),
				),
			),
			result: &target{},
			output: &target{
				Foo: "bar",
				Baz: 32,
				Sub: &sub{Foo: 13},
			},
		}, {
			name: "non-empty document into empty struct with non-empty subdocument",
			input: NewDocument(
				EC.String("foo", "bar"),
				EC.Int32("baz", 32),
				EC.SubDocumentFromElements(
					"sub",
					EC.Int32("foo", 13),
				),
			),
			result: &target{
				Foo: "overwrite",
				Bar: 13,
				Sub: &sub{Foo: 12, Bar: 10},
			},
			output: &target{
				Foo: "bar",
				Bar: 13,
				Baz: 32,
				Sub: &sub{Foo: 13, Bar: 10},
			},
		}}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				var buffer bytes.Buffer
				test.input.WriteTo(&buffer)
				actual, err := codec.Decode(&buffer, test.result)
				require.NoError(t, err)

				require.Equal(t, test.output, test.result)
				require.Equal(t, test.result, actual)
			})
		}
	})
}
