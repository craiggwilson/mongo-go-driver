package bson_test

import (
	"strconv"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson"
)

var result int

func BenchmarkWideDoc(b *testing.B) {
	var d bson.Doc
	for n := 0; n < b.N; n++ {
		d = buildWideDoc()
	}
	result = len(d)
}

func BenchmarkWideDoc2(b *testing.B) {
	var d bson.Doc2
	for n := 0; n < b.N; n++ {
		d = buildWideDoc2()
	}
	result = len(d)
}

func BenchmarkWideD(b *testing.B) {
	var d bson.D
	for n := 0; n < b.N; n++ {
		d = buildWideD()
	}
	result = len(d)
}

func buildWideDoc() bson.Doc {
	d := bson.Doc{}
	for i := 0; i < 100; i++ {
		d = append(d, bson.Elem{strconv.Itoa(i), bson.Int32(int32(i))})
	}
	return d
}

func buildWideDoc2() bson.Doc2 {
	d := bson.Doc2{}
	for i := 0; i < 100; i++ {
		d = append(d, bson.Elem2{strconv.Itoa(i), bson.Int32(int32(i))})
	}
	return d
}

func buildWideD() bson.D {
	d := bson.D{}
	for i := 0; i < 100; i++ {
		d = append(d, bson.E{strconv.Itoa(i), i})
	}
	return d
}

func BenchmarkDeepDoc(b *testing.B) {
	var d bson.Doc
	for n := 0; n < b.N; n++ {
		d = buildDeepDoc(100)
	}
	result = len(d)
}

func BenchmarkDeepDoc2(b *testing.B) {
	var d bson.Doc2
	for n := 0; n < b.N; n++ {
		d = buildDeepDoc2(100)
	}
	result = len(d)
}

func BenchmarkDeepD(b *testing.B) {
	var d bson.D
	for n := 0; n < b.N; n++ {
		d = buildDeepD(100)
	}
	result = len(d)
}

func buildDeepDoc(level int) bson.Doc {
	d := bson.Doc{}
	if level > 0 {
		d = append(d, bson.Elem{"i", bson.Document(buildDeepDoc(level - 1))})
	}
	return d
}

func buildDeepDoc2(level int) bson.Doc2 {
	d := bson.Doc2{}
	if level > 0 {
		d = append(d, bson.Elem2{"i", buildDeepDoc2(level - 1)})
	}
	return d
}

func buildDeepD(level int) bson.D {
	d := bson.D{}
	if level > 0 {
		d = append(d, bson.E{"i", buildDeepD(level - 1)})
	}
	return d
}
