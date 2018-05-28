package bench

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson2"
	mgo "gopkg.in/mgo.v2/bson"
)

func benchmarkWriteBsonDocument(input *bson.Document, b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		_, err = bson.Marshal(input)
	}
	benchError = err
}

func benchmarkWriteBson2D(input bson2.D, b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		_, err = bson2.Marshal(input)
	}
	benchError = err
}

func benchmarkWriteMgoD(input mgo.D, b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		_, err = mgo.Marshal(input)
	}
	benchError = err
}

func BenchmarkWriteSmall_Bson_Document(b *testing.B) {
	input := bson.NewDocument(
		bson.EC.Int32("a", 1),
		bson.EC.SubDocumentFromElements("x",
			bson.EC.String("a", "b"),
		),
	)

	benchmarkWriteBsonDocument(input, b)
}

func BenchmarkWriteSmall_Bson2_D(b *testing.B) {
	input := bson2.D{{Name: "a", Value: bson2.D{{Name: "a", Value: "b"}}}}
	benchmarkWriteBson2D(input, b)
}

func BenchmarkWriteSmall_Mgo_D(b *testing.B) {
	input := mgo.D{{Name: "a", Value: mgo.D{{Name: "a", Value: "b"}}}}
	benchmarkWriteMgoD(input, b)
}
