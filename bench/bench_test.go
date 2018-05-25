package bench

import (
	"bytes"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson2"
	mgo "gopkg.in/mgo.v2/bson"
)

func benchmarkBson2D(input []byte, b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		target := &bson2.D{}
		bson2.Unmarshal(r, target)
	}
}

func benchmarkBson2RawD(input []byte, b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(input)
		target := &bson2.RawD{}
		bson2.Unmarshal(r, target)
	}
}

func benchmarkBsonDocument(input []byte, b *testing.B) {
	for i := 0; i < b.N; i++ {
		target := bson.NewDocument()
		bson.Unmarshal(input, &target)
	}
}

func benchmarkMgoD(input []byte, b *testing.B) {
	for i := 0; i < b.N; i++ {
		target := &mgo.D{}
		mgo.Unmarshal(input, target)
	}
}

func benchmarkMgoRawD(input []byte, b *testing.B) {
	for i := 0; i < b.N; i++ {
		target := &mgo.RawD{}
		mgo.Unmarshal(input, target)
	}
}

func BenchmarkSmallBsonDocument(b *testing.B) { benchmarkBsonDocument(small, b) }
func BenchmarkSmallBson2D(b *testing.B)       { benchmarkBson2D(small, b) }
func BenchmarkSmallBson2RawD(b *testing.B)    { benchmarkBson2RawD(small, b) }
func BenchmarkSmallMgoD(b *testing.B)         { benchmarkMgoD(small, b) }
func BenchmarkSmallMgoRawD(b *testing.B)      { benchmarkMgoRawD(small, b) }

func BenchmarkSmall2BsonDocument(b *testing.B) { benchmarkBsonDocument(small2, b) }
func BenchmarkSmall2Bson2D(b *testing.B)       { benchmarkBson2D(small2, b) }
func BenchmarkSmall2Bson2RawD(b *testing.B)    { benchmarkBson2RawD(small2, b) }
func BenchmarkSmall2MgoD(b *testing.B)         { benchmarkMgoD(small2, b) }
func BenchmarkSmall2MgoRawD(b *testing.B)      { benchmarkMgoRawD(small2, b) }

func BenchmarkLargeFlatBsonDocument(b *testing.B) { benchmarkBsonDocument(largeFlat, b) }
func BenchmarkLargeFlatBson2D(b *testing.B)       { benchmarkBson2D(largeFlat, b) }
func BenchmarkLargeFlatBson2RawD(b *testing.B)    { benchmarkBson2RawD(largeFlat, b) }
func BenchmarkLargeFlatMgoD(b *testing.B)         { benchmarkMgoD(largeFlat, b) }
func BenchmarkLargeFlatMgoRawD(b *testing.B)      { benchmarkMgoRawD(largeFlat, b) }

func BenchmarkLargeDeepBsonDocument(b *testing.B) { benchmarkBsonDocument(largeDeep, b) }
func BenchmarkLargeDeepBson2D(b *testing.B)       { benchmarkBson2D(largeDeep, b) }
func BenchmarkLargeDeepBson2RawD(b *testing.B)    { benchmarkBson2RawD(largeDeep, b) }
func BenchmarkLargeDeepMgoD(b *testing.B)         { benchmarkMgoD(largeDeep, b) }
func BenchmarkLargeDeepMgoRawD(b *testing.B)      { benchmarkMgoRawD(largeDeep, b) }