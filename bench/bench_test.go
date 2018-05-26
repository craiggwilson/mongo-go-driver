package bench

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson2"
	mgo "gopkg.in/mgo.v2/bson"
)

func benchmarkBsonDocument(input []byte, b *testing.B) {
	for i := 0; i < b.N; i++ {
		target := bson.NewDocument()
		bson.Unmarshal(input, target)
	}
}

func benchmarkBson2D(input []byte, b *testing.B) {
	for i := 0; i < b.N; i++ {
		target := &bson2.D{}
		bson2.Unmarshal(input, target)
	}
}

func benchmarkBson2M(input []byte, b *testing.B) {
	for i := 0; i < b.N; i++ {
		target := &bson2.M{}
		bson2.Unmarshal(input, target)
	}
}

func benchmarkBson2RawD(input []byte, b *testing.B) {
	for i := 0; i < b.N; i++ {
		target := &bson2.RawD{}
		bson2.Unmarshal(input, target)
	}
}

func benchmarkBson2Raw(input []byte, b *testing.B) {
	for i := 0; i < b.N; i++ {
		target := &bson2.Raw{}
		bson2.Unmarshal(input, target)
	}
}

func benchmarkMgoD(input []byte, b *testing.B) {
	for i := 0; i < b.N; i++ {
		target := &mgo.D{}
		mgo.Unmarshal(input, target)
	}
}

func benchmarkMgoM(input []byte, b *testing.B) {
	for i := 0; i < b.N; i++ {
		target := &mgo.M{}
		mgo.Unmarshal(input, target)
	}
}

func benchmarkMgoRawD(input []byte, b *testing.B) {
	for i := 0; i < b.N; i++ {
		target := &mgo.RawD{}
		mgo.Unmarshal(input, target)
	}
}

func BenchmarkSmall_Bson_Document(b *testing.B) { benchmarkBsonDocument(small, b) }
func BenchmarkSmall_Bson2_D(b *testing.B)       { benchmarkBson2D(small, b) }
func BenchmarkSmall_Bson2_M(b *testing.B)       { benchmarkBson2M(small, b) }
func BenchmarkSmall_Bson2_RawD(b *testing.B)    { benchmarkBson2RawD(small, b) }
func BenchmarkSmall_Bson2_Raw(b *testing.B)     { benchmarkBson2Raw(small, b) }
func BenchmarkSmall_Mgo_D(b *testing.B)         { benchmarkMgoD(small, b) }
func BenchmarkSmall_Mgo_M(b *testing.B)         { benchmarkMgoM(small, b) }
func BenchmarkSmall_Mgo_RawD(b *testing.B)      { benchmarkMgoRawD(small, b) }

func BenchmarkSmall2_Bson_Document(b *testing.B) { benchmarkBsonDocument(small2, b) }
func BenchmarkSmall2_Bson2_D(b *testing.B)       { benchmarkBson2D(small2, b) }
func BenchmarkSmall2_Bson2_M(b *testing.B)       { benchmarkBson2M(small2, b) }
func BenchmarkSmall2_Bson2_RawD(b *testing.B)    { benchmarkBson2RawD(small2, b) }
func BenchmarkSmall2_Bson2_Raw(b *testing.B)     { benchmarkBson2Raw(small2, b) }
func BenchmarkSmall2_Mgo_D(b *testing.B)         { benchmarkMgoD(small2, b) }
func BenchmarkSmall2_Mgo_M(b *testing.B)         { benchmarkMgoM(small2, b) }
func BenchmarkSmall2_Mgo_RawD(b *testing.B)      { benchmarkMgoRawD(small2, b) }

func BenchmarkLargeFlat_Bson_Document(b *testing.B) { benchmarkBsonDocument(largeFlat, b) }
func BenchmarkLargeFlat_Bson2_D(b *testing.B)       { benchmarkBson2D(largeFlat, b) }
func BenchmarkLargeFlat_Bson2_M(b *testing.B)       { benchmarkBson2M(largeFlat, b) }
func BenchmarkLargeFlat_Bson2_RawD(b *testing.B)    { benchmarkBson2RawD(largeFlat, b) }
func BenchmarkLargeFlat_Bson2_Raw(b *testing.B)     { benchmarkBson2Raw(largeFlat, b) }
func BenchmarkLargeFlat_Mgo_D(b *testing.B)         { benchmarkMgoD(largeFlat, b) }
func BenchmarkLargeFlat_Mgo_M(b *testing.B)         { benchmarkMgoM(largeFlat, b) }
func BenchmarkLargeFlat_Mgo_RawD(b *testing.B)      { benchmarkMgoRawD(largeFlat, b) }

func BenchmarkLargeDeep_Bson_Document(b *testing.B) { benchmarkBsonDocument(largeDeep, b) }
func BenchmarkLargeDeep_Bson_2D(b *testing.B)       { benchmarkBson2D(largeDeep, b) }
func BenchmarkLargeDeep_Bson_2M(b *testing.B)       { benchmarkBson2M(largeDeep, b) }
func BenchmarkLargeDeep_Bson_2RawD(b *testing.B)    { benchmarkBson2RawD(largeDeep, b) }
func BenchmarkLargeDeep_Bson_2Raw(b *testing.B)     { benchmarkBson2Raw(largeDeep, b) }
func BenchmarkLargeDeep_Mgo_D(b *testing.B)         { benchmarkMgoD(largeDeep, b) }
func BenchmarkLargeDeep_Mgo_M(b *testing.B)         { benchmarkMgoM(largeDeep, b) }
func BenchmarkLargeDeep_Mgo_RawD(b *testing.B)      { benchmarkMgoRawD(largeDeep, b) }
