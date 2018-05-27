package bench

import (
	"reflect"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson2"
	mgo "gopkg.in/mgo.v2/bson"
)

var benchError error

func benchmarkReadBsonDocument(input []byte, b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		target := bson.NewDocument()
		err = bson.Unmarshal(input, target)
	}
	benchError = err
}

func benchmarkReadBsonStruct(input []byte, t reflect.Type, b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		target := reflect.New(t).Interface()
		err = bson.Unmarshal(input, target)
	}
	benchError = err
}

func benchmarkReadBson2Document(input []byte, b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		target := &bson2.Document{}
		err = bson2.Unmarshal(input, target)
	}
	benchError = err
}

func benchmarkReadBson2Struct(input []byte, t reflect.Type, b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		target := reflect.New(t).Interface()
		err = bson2.Unmarshal(input, target)
	}
	benchError = err
}

func benchmarkReadBson2D(input []byte, b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		target := &bson2.D{}
		err = bson2.Unmarshal(input, target)
	}
	benchError = err
}

func benchmarkBson2M(input []byte, b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		target := &bson2.M{}
		err = bson2.Unmarshal(input, target)
	}
	benchError = err
}

func benchmarkReadBson2RawD(input []byte, b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		target := &bson2.RawD{}
		err = bson2.Unmarshal(input, target)
	}
	benchError = err
}

func benchmarkBson2Raw(input []byte, b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		target := &bson2.Raw{}
		err = bson2.Unmarshal(input, target)
	}
	benchError = err
}

func benchmarkReadMgoD(input []byte, b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		target := &mgo.D{}
		err = mgo.Unmarshal(input, target)
	}
	benchError = err
}

func benchmarkReadMgoM(input []byte, b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		target := &mgo.M{}
		err = mgo.Unmarshal(input, target)
	}
	benchError = err
}

func benchmarkReadMgoRawD(input []byte, b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		target := &mgo.RawD{}
		err = mgo.Unmarshal(input, target)
	}
	benchError = err
}

func BenchmarkWarmup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = small
		_ = small2
		_ = largeFlat
		_ = largeDeep
	}
}

func BenchmarkReadSmall_Bson_Document(b *testing.B) { benchmarkReadBsonDocument(small, b) }
func BenchmarkReadSmall_Bson_Struct(b *testing.B) {
	benchmarkReadBsonStruct(small, reflect.TypeOf(new(smallStruct)), b)
}
func BenchmarkReadSmall_Bson2_Document(b *testing.B) { benchmarkReadBson2Document(small, b) }
func BenchmarkReadSmall_Bson2_Struct(b *testing.B) {
	benchmarkReadBson2Struct(small, reflect.TypeOf(new(smallStruct)), b)
}
func BenchmarkReadSmall_Bson2_D(b *testing.B)    { benchmarkReadBson2D(small, b) }
func BenchmarkReadSmall_Bson2_M(b *testing.B)    { benchmarkBson2M(small, b) }
func BenchmarkReadSmall_Bson2_RawD(b *testing.B) { benchmarkReadBson2RawD(small, b) }
func BenchmarkReadSmall_Bson2_Raw(b *testing.B)  { benchmarkBson2Raw(small, b) }
func BenchmarkReadSmall_Mgo_D(b *testing.B)      { benchmarkReadMgoD(small, b) }
func BenchmarkReadSmall_Mgo_M(b *testing.B)      { benchmarkReadMgoM(small, b) }
func BenchmarkReadSmall_Mgo_RawD(b *testing.B)   { benchmarkReadMgoRawD(small, b) }

func BenchmarkReadSmall2_Bson_Document(b *testing.B) { benchmarkReadBsonDocument(small2, b) }
func BenchmarkReadSmall2_Bson_Struct(b *testing.B) {
	benchmarkReadBsonStruct(small2, reflect.TypeOf(new(small2Struct)), b)
}
func BenchmarkReadSmall2_Bson2_Document(b *testing.B) { benchmarkReadBson2Document(small2, b) }
func BenchmarkReadSmall2_Bson2_Struct(b *testing.B) {
	benchmarkReadBson2Struct(small2, reflect.TypeOf(new(small2Struct)), b)
}
func BenchmarkReadSmall2_Bson2_D(b *testing.B)    { benchmarkReadBson2D(small2, b) }
func BenchmarkReadSmall2_Bson2_M(b *testing.B)    { benchmarkBson2M(small2, b) }
func BenchmarkReadSmall2_Bson2_RawD(b *testing.B) { benchmarkReadBson2RawD(small2, b) }
func BenchmarkReadSmall2_Bson2_Raw(b *testing.B)  { benchmarkBson2Raw(small2, b) }
func BenchmarkReadSmall2_Mgo_D(b *testing.B)      { benchmarkReadMgoD(small2, b) }
func BenchmarkReadSmall2_Mgo_M(b *testing.B)      { benchmarkReadMgoM(small2, b) }
func BenchmarkReadSmall2_Mgo_RawD(b *testing.B)   { benchmarkReadMgoRawD(small2, b) }

func BenchmarkReadLargeFlat_Bson_Document(b *testing.B)  { benchmarkReadBsonDocument(largeFlat, b) }
func BenchmarkReadLargeFlat_Bson2_Document(b *testing.B) { benchmarkReadBson2Document(largeFlat, b) }
func BenchmarkReadLargeFlat_Bson2_D(b *testing.B)        { benchmarkReadBson2D(largeFlat, b) }
func BenchmarkReadLargeFlat_Bson2_M(b *testing.B)        { benchmarkBson2M(largeFlat, b) }
func BenchmarkReadLargeFlat_Bson2_RawD(b *testing.B)     { benchmarkReadBson2RawD(largeFlat, b) }
func BenchmarkReadLargeFlat_Bson2_Raw(b *testing.B)      { benchmarkBson2Raw(largeFlat, b) }
func BenchmarkReadLargeFlat_Mgo_D(b *testing.B)          { benchmarkReadMgoD(largeFlat, b) }
func BenchmarkReadLargeFlat_Mgo_M(b *testing.B)          { benchmarkReadMgoM(largeFlat, b) }
func BenchmarkReadLargeFlat_Mgo_RawD(b *testing.B)       { benchmarkReadMgoRawD(largeFlat, b) }

func BenchmarkReadLargeDeep_Bson_Document(b *testing.B)  { benchmarkReadBsonDocument(largeDeep, b) }
func BenchmarkReadLargeDeep_Bson2_Document(b *testing.B) { benchmarkReadBson2Document(largeDeep, b) }
func BenchmarkReadLargeDeep_Bson_2D(b *testing.B)        { benchmarkReadBson2D(largeDeep, b) }
func BenchmarkReadLargeDeep_Bson_2M(b *testing.B)        { benchmarkBson2M(largeDeep, b) }
func BenchmarkReadLargeDeep_Bson_2RawD(b *testing.B)     { benchmarkReadBson2RawD(largeDeep, b) }
func BenchmarkReadLargeDeep_Bson_2Raw(b *testing.B)      { benchmarkBson2Raw(largeDeep, b) }
func BenchmarkReadLargeDeep_Mgo_D(b *testing.B)          { benchmarkReadMgoD(largeDeep, b) }
func BenchmarkReadLargeDeep_Mgo_M(b *testing.B)          { benchmarkReadMgoM(largeDeep, b) }
func BenchmarkReadLargeDeep_Mgo_RawD(b *testing.B)       { benchmarkReadMgoRawD(largeDeep, b) }
