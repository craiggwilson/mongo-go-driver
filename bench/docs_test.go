package bench

import (
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
)

var small, _ = bson.NewDocument(
	bson.EC.Int32("a", 1),
	bson.EC.SubDocumentFromElements("x",
		bson.EC.String("a", "b"),
	),
).MarshalBSON()

var small2, _ = bson.NewDocument(
	bson.EC.SubDocumentFromElements("driver",
		bson.EC.String("name", "mongo-go-driver"),
		bson.EC.String("version", "234234"),
	),
	bson.EC.SubDocumentFromElements("os",
		bson.EC.String("type", "darwin"),
		bson.EC.String("architecture", "amd64"),
	),
	bson.EC.String("platform", "go1.9.2"),
).MarshalBSON()

var largeFlat, _ = buildLargeFlatDocument().MarshalBSON()

var largeDeep, _ = buildLargeDeepDocument().MarshalBSON()

func buildLargeFlatDocument() *bson.Document {
	doc := bson.NewDocument()
	for i := 0; i < 2000; i++ {
		doc.Append(bson.EC.Int32(fmt.Sprintf("a%d", i), int32(i)))
	}

	return doc
}

func buildLargeDeepDocument() *bson.Document {
	var subdoc func(depth int) *bson.Document
	subdoc = func(depth int) *bson.Document {
		doc := buildLargeFlatDocument()

		if depth < 50 {
			doc.Append(bson.EC.SubDocument("b", subdoc(depth+1)))
		}

		return doc
	}

	d := subdoc(0)
	return d
}

// func TestPrint(t *testing.T) {
// 	target := &bson2.D{}
// 	bson2.Unmarshal(largeDeep, target)

// 	fmt.Printf("%# v", target)
// 	t.Fail()
// }

// func TestPrint(t *testing.T) {
// 	d := buildLargeDeepDocument()
// 	fmt.Println(d.String())
// 	t.Fail()
// }
