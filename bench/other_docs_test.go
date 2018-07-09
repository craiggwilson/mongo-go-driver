package bench

import (
	"reflect"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson"
)

func BenchmarkOtherReadSmall_Bson_Struct(b *testing.B) {
	benchmarkReadBsonStruct(aSmallStructBytes, reflect.TypeOf(new(SmallStruct)), b)
}
func BenchmarkOtherReadSmall_Bson2_Struct(b *testing.B) {
	benchmarkReadBson2Struct(aSmallStructBytes, reflect.TypeOf(new(SmallStruct)), b)
}
func BenchmarkOtherReadSmall_Mgo_Struct(b *testing.B) {
	benchmarkReadMgoStruct(aSmallStructBytes, reflect.TypeOf(new(SmallStruct)), b)
}
func BenchmarkOtherReadSmallNested_Bson_Struct(b *testing.B) {
	benchmarkReadBsonStruct(aSmallNestedStructBytes, reflect.TypeOf(new(SmallStructDepth9)), b)
}
func BenchmarkOtherReadSmallNested_Bson2_Struct(b *testing.B) {
	benchmarkReadBson2Struct(aSmallNestedStructBytes, reflect.TypeOf(new(SmallStructDepth9)), b)
}
func BenchmarkOtherReadSmallNested_Mgo_Struct(b *testing.B) {
	benchmarkReadMgoStruct(aSmallNestedStructBytes, reflect.TypeOf(new(SmallStructDepth9)), b)
}
func BenchmarkOtherReadLarge_Bson_Struct(b *testing.B) {
	benchmarkReadBsonStruct(aLargerStructBytes, reflect.TypeOf(new(LargerStruct)), b)
}
func BenchmarkOtherReadLarge_Bson2_Struct(b *testing.B) {
	benchmarkReadBson2Struct(aLargerStructBytes, reflect.TypeOf(new(LargerStruct)), b)
}
func BenchmarkOtherReadLarge_Mgo_Struct(b *testing.B) {
	benchmarkReadMgoStruct(aLargerStructBytes, reflect.TypeOf(new(LargerStruct)), b)
}
func BenchmarkOtherReadLargeNested_Bson_Struct(b *testing.B) {
	benchmarkReadBsonStruct(aLargerNestedStructBytes, reflect.TypeOf(new(LargerStructDepth9)), b)
}
func BenchmarkOtherReadLargeNested_Bson2_Struct(b *testing.B) {
	benchmarkReadBson2Struct(aLargerNestedStructBytes, reflect.TypeOf(new(LargerStructDepth9)), b)
}
func BenchmarkOtherReadLargeNested_Mgo_Struct(b *testing.B) {
	benchmarkReadMgoStruct(aLargerNestedStructBytes, reflect.TypeOf(new(LargerStructDepth9)), b)
}

var textMediumLength = "The quick brown fox jumps over the lazy dog. Then it did it again. Keep typing a " +
	"little longer so we get a string that's somewhere between a short string (like a hostname perhaps) and a long " +
	"text blob"

// zero values are fine for everything but strings

var aSmallStructBytes, _ = bson.Marshal(ASmallStruct)
var aLargerStructBytes, _ = bson.Marshal(ALargerStruct)
var aSmallNestedStructBytes, _ = bson.Marshal(ASmallNestedStruct)
var aLargerNestedStructBytes, _ = bson.Marshal(ALargerNestedStruct)

var ASmallStruct = SmallStruct{
	SomeText: textMediumLength,
}

var ALargerStruct = LargerStruct{
	SomeText1:  textMediumLength,
	SomeText2:  textMediumLength,
	SomeText3:  textMediumLength,
	SomeText4:  textMediumLength,
	SomeText5:  textMediumLength,
	SomeText6:  textMediumLength,
	SomeText7:  textMediumLength,
	SomeText8:  textMediumLength,
	SomeText9:  textMediumLength,
	SomeText10: textMediumLength,
}

var ASmallNestedStruct = SmallStructDepth9{
	Other: ASmallStruct,
	// nested structs will have empty strings, meh
}

var ALargerNestedStruct = LargerStructDepth9{
	SomeText1:  textMediumLength,
	SomeText2:  textMediumLength,
	SomeText3:  textMediumLength,
	SomeText4:  textMediumLength,
	SomeText5:  textMediumLength,
	SomeText6:  textMediumLength,
	SomeText7:  textMediumLength,
	SomeText8:  textMediumLength,
	SomeText9:  textMediumLength,
	SomeText10: textMediumLength,
	// nested structs will have empty strings, meh
}

type SmallStruct struct {
	ANumber  int32
	ATime    int32
	SomeText string
}

type LargerStruct struct {
	ANumber1   int32
	ATime1     int32
	SomeText1  string
	ANumber2   int32
	ATime2     int32
	SomeText2  string
	ABool      int32
	ATime3     int32
	SomeText3  string
	ANumber4   int32
	ATime4     int32
	SomeText4  string
	ANumber5   int32
	ATime5     int32
	SomeText5  string
	ANumber6   int32
	ATime6     int32
	SomeText6  string
	AFloat     int32
	ATime7     int32
	SomeText7  string
	ANumber    int32
	ATime8     int32
	SomeText8  string
	ANumber9   int32
	ATime9     int32
	SomeText9  string
	AByte      int32
	ATime10    int32
	SomeText10 string
}

type SmallStructDepth9 struct {
	Other  SmallStruct
	Nested smallNestedStruct2
}

type smallNestedStruct2 struct {
	Other  SmallStruct
	Nested smallNestedStruct3
}

type smallNestedStruct3 struct {
	Other  SmallStruct
	Nested smallNestedStruct4
}

type smallNestedStruct4 struct {
	Other  SmallStruct
	Nested smallNestedStruct5
}

type smallNestedStruct5 struct {
	Other  SmallStruct
	Nested smallNestedStruct6
}

type smallNestedStruct6 struct {
	Other  SmallStruct
	Nested smallNestedStruct7
}

type smallNestedStruct7 struct {
	Other  SmallStruct
	Nested smallNestedStruct8
}

type smallNestedStruct8 struct {
	Other  SmallStruct
	Nested smallNestedStruct9
}

type smallNestedStruct9 struct {
	Other  SmallStruct
	Nested SmallStruct
}

type LargerStructDepth9 struct {
	ANumber1   int32
	ATime1     int32
	SomeText1  string
	ANumber2   int32
	ATime2     int32
	SomeText2  string
	ABool      int32
	ATime3     int32
	SomeText3  string
	ANumber4   int32
	ATime4     int32
	SomeText4  string
	ANumber5   int32
	ATime5     int32
	SomeText5  string
	ANumber6   int32
	ATime6     int32
	SomeText6  string
	AFloat     int32
	ATime7     int32
	SomeText7  string
	ANumber    int32
	ATime8     int32
	SomeText8  string
	ANumber9   int32
	ATime9     int32
	SomeText9  string
	AByte      int32
	ATime10    int32
	SomeText10 string
	Nested     largerNestedStruct2
}

type largerNestedStruct2 struct {
	Other  LargerStruct
	Nested largerNestedStruct3
}

type largerNestedStruct3 struct {
	Other  LargerStruct
	Nested largerNestedStruct4
}

type largerNestedStruct4 struct {
	Other  LargerStruct
	Nested largerNestedStruct5
}

type largerNestedStruct5 struct {
	Other  LargerStruct
	Nested largerNestedStruct6
}

type largerNestedStruct6 struct {
	Other  LargerStruct
	Nested largerNestedStruct7
}

type largerNestedStruct7 struct {
	Other  LargerStruct
	Nested largerNestedStruct8
}

type largerNestedStruct8 struct {
	Other  LargerStruct
	Nested largerNestedStruct9
}

type largerNestedStruct9 struct {
	Other  LargerStruct
	Nested LargerStruct
}
