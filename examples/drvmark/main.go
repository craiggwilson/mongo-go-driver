package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"sync"

	"runtime"

	"github.com/10gen/mongo-go-driver/cluster"
	"github.com/10gen/mongo-go-driver/conn"
	"github.com/10gen/mongo-go-driver/connstring"
	"github.com/10gen/mongo-go-driver/msg"
	"github.com/10gen/mongo-go-driver/ops"
	"github.com/10gen/mongo-go-driver/readpref"
)

var batch = flag.Bool("batch", true, "whether to use batches for insert")
var concurrency = flag.Int("conc", 4, "level of concurrency")
var n = flag.Int("n", 1000000, "number of documents")
var uri = flag.String("uri", "mongodb://localhost:27017", "mongodb uri")

func main() {

	flag.Parse()

	runtime.GOMAXPROCS(*concurrency)

	connStr, err := connstring.Parse(*uri)
	if err != nil {
		log.Fatalf("error parsing uri: %v", err)
	}

	var endpoints []conn.Endpoint
	for _, host := range connStr.Hosts {
		endpoints = append(endpoints, conn.Endpoint(host))
	}

	c, err := cluster.New(
		cluster.WithSeedList(endpoints...),
	)
	if err != nil {
		log.Fatalf("error creating cluster: %v", err)
	}

	var totalLinearElapsed time.Duration
	var totalParallelElapsed time.Duration

	var stTime time.Duration
	var mtTime time.Duration

	ctx := context.Background()
	sections := []string{"create", "update", "read"}
	for _, section := range sections {
		log.Println("Starting section (Linear): ", section)
		stTime = singleWorker(ctx, c, section, 0)
		totalLinearElapsed += stTime
		log.Printf("Ending section (Linear): %s, elapsed: %s\n", section, stTime)
	}
	log.Printf("Total Elapsed (Linear): %s\n", totalLinearElapsed)

	for _, section := range sections {
		log.Println("Starting section (Parallel): ", section)
		mtTime = multiWorker(ctx, c, section, 0)
		totalParallelElapsed += mtTime
		log.Printf("Ending section (Parallel): %s, elapsed: %s\n", section, mtTime)
	}
	log.Printf("Total Elapsed (Parallel): %s\n", totalParallelElapsed)
}

func singleWorker(ctx context.Context, c *cluster.Cluster, section string, idx int) time.Duration {
	ns := ops.Namespace{
		DB:         "dvrmark-go",
		Collection: fmt.Sprintf("records_%d", idx),
	}

	var totalElapsed time.Duration

	if section == "create" || section == "all" {
		totalElapsed += create(ctx, c, ns, idx)
	}

	if section == "read" || section == "all" {
		totalElapsed += read(ctx, c, ns, idx)
	}

	return totalElapsed
}

func multiWorker(ctx context.Context, c *cluster.Cluster, section string, idx int) time.Duration {
	start := time.Now()
	rem := *n % *concurrency
	chunk := *n / *concurrency

	var wg sync.WaitGroup
	for i := 1; i < *concurrency; i++ {
		wg.Add(1)
		if i == *concurrency {
			chunk = chunk + rem
		}

		go func() {
			singleWorker(ctx, c, section, idx)
			wg.Done()
		}()
	}
	wg.Wait()
	return time.Since(start)
}

func insert(ctx context.Context, c *cluster.Cluster, ns ops.Namespace, docs []bson.M, idx int) {
	insertCommand := bson.D{
		{"insert", ns.Collection},
		{"documents", docs},
	}
	request := msg.NewCommand(
		msg.NextRequestID(),
		ns.DB,
		false,
		insertCommand,
	)

	s, err := c.SelectServer(ctx, cluster.WriteSelector())
	if err != nil {
		log.Fatalf("t%d -> unable to select server for insert: %v", idx, err)
	}

	connection, err := s.Connection(ctx)
	if err != nil {
		log.Fatalf("t%d -> unable to get connection for insert: %v", idx, err)
	}
	defer connection.Close()

	err = conn.ExecuteCommand(ctx, connection, request, &bson.D{})
	if err != nil {
		log.Fatalf("t%d -> unable to execute insert: %v", idx, err)
	}
}

func insertBulk(ctx context.Context, c *cluster.Cluster, ns ops.Namespace, idx int, batchSize int) {
	docs := make([]bson.M, 0, batchSize)
	for i := 1; i < *n; i++ {
		docs = append(docs, generateDoc(idx, i))
		if i%batchSize == 0 {
			insert(ctx, c, ns, docs, idx)
		}
		docs = make([]bson.M, 0, batchSize)
	}

	if len(docs) > 0 {
		insert(ctx, c, ns, docs, idx)
	}
}

func create(ctx context.Context, c *cluster.Cluster, ns ops.Namespace, idx int) time.Duration {
	start := time.Now()
	drop(ctx, c, ns, idx)
	insertBulk(ctx, c, ns, idx, 1000)
	return time.Since(start)
}

func drop(ctx context.Context, c *cluster.Cluster, ns ops.Namespace, idx int) {
	s, err := c.SelectServer(ctx, cluster.WriteSelector())
	if err != nil {
		log.Fatalf("t%d -> unable to drop ns: %v", idx, err)
	}
	connection, err := s.Connection(ctx)
	if err != nil {
		log.Fatalf("t%d -> unable to get connection for drop: %v", idx, err)
	}
	defer connection.Close()

	err = conn.ExecuteCommand(
		ctx,
		connection,
		msg.NewCommand(
			msg.NextRequestID(),
			ns.DB,
			false,
			bson.D{{"drop", ns.Collection}},
		),
		&bson.D{},
	)
	if err != nil && !strings.HasSuffix(err.Error(), "ns not found") {
		log.Fatalf("t%d -> failed dropping ns: %v", idx, err)
	}
}

func read(ctx context.Context, c *cluster.Cluster, ns ops.Namespace, idx int) time.Duration {
	start := time.Now()
	rp := readpref.Primary()
	svr, err := c.SelectServer(ctx, cluster.ReadPrefSelector(rp))
	if err != nil {
		log.Fatalf("t%d -> error selecting a server: %v", idx, err)
	}

	cur, err := ops.Aggregate(ctx, &ops.SelectedServer{svr, rp}, ns, []bson.D{}, ops.AggregationOptions{})
	if err != nil {
		log.Fatalf("t%d -> error executing aggregate: %v", idx, err)
	}

	var result interface{}
	var sum float64

	for cur.Next(ctx, &result) {
		doc := result.(bson.M)
		subArray := doc["arr"].([]interface{})
		for _, v := range subArray {
			doc := v.(bson.M)
			sum = sum + doc["subval2"].(float64)
		}
	}
	if cur.Err() != nil {
		log.Fatalf("t%d -> error iterating cursor: %v", idx, cur.Err())
	}
	cur.Close(ctx)

	return time.Since(start)
}

func generateDoc(idx int, docIdx int) bson.M {
	const topFields = 20
	const arrSize = 20
	const arrObjSize = 10
	const fieldPrefix = "val"

	id := fmt.Sprintf("%d-%d", docIdx%256, docIdx)

	doc := bson.M{"_id": id}

	for i := 0; i < topFields; i++ {
		tp := i % 4
		fieldName := fmt.Sprintf("%s%d", fieldPrefix, i)
		switch tp {
		case 0:
			v := "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
			doc[fieldName] = v
		case 1:
			v := time.Unix(int64(i*docIdx/1000), 0)
			doc[fieldName] = v
		case 2:
			v := math.Pi * float64(i)
			doc[fieldName] = v
		case 3:
			v := int64(docIdx + i)
			doc[fieldName] = v
		}
	}

	var arrObject []interface{}

	for i := 0; i < arrSize; i++ {
		subDoc := bson.M{"name": "subRec"}
		for j := 0; j < arrObjSize; j++ {
			tp := j % 4
			fieldName := fmt.Sprintf("subval%d", tp)
			switch tp {
			case 0:
				v := "Nunc finibus pretium dignissim. Aenean ut nisi finibus."
				subDoc[fieldName] = v
			case 1:
				v := time.Unix(int64(i*docIdx/1000), 0)
				subDoc[fieldName] = v
			case 2:
				v := math.Pi * float64(i)
				subDoc[fieldName] = v
			case 3:
				v := int64(docIdx + i)
				subDoc[fieldName] = v
			}
		}
		arrObject = append(arrObject, subDoc)
	}

	doc["arr"] = arrObject

	return doc
}
