package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type Args struct {
	mongoURI   *string
	db         *string
	collection *string
	key        *string
	blockSize  *int
}

var args Args

func parseArgs() {
	args = Args{
		mongoURI:   flag.String("mongoURI", "", "MongoDB URI"),
		db:         flag.String("db", "", "MongoDB database"),
		collection: flag.String("collection", "", "MongoDB collection"),
		key:        flag.String("key", "", "DynamoDB partition key that will be used as _id"),
		blockSize:  flag.Int("blockSize", 100, "How many items to insert on each request to Mongo"),
	}

	flag.Parse()
}

func initScanner(r io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(r)
	const maxDDBDocSize = 400 * 1024
	buffer := make([]byte, maxDDBDocSize)
	scanner.Buffer(buffer, maxDDBDocSize)

	return scanner
}

func buildDoc(raw []byte, key string) interface{} {
	ddbJson := Unmarshal(raw)
	jsonMap := ddbJson.ToMap() // We start always from a map
	doc := jsonMap.ToBSON(key)

	return doc
}

func transmit(r io.Reader, key string, c chan interface{}, quit chan struct{}) {
	scanner := initScanner(r)

	for scanner.Scan() {
		doc := buildDoc(scanner.Bytes(), key)
		c <- &doc
	}

	quit <- struct{}{}
}

func process(c chan interface{}, quit chan struct{}, loader JSONLoader, db string, collection string, key string, blockSize int) {
	var size int
	docs := make([]interface{}, 0, blockSize)

	handleDocs := func() {
		loader.BulkLoad(db, collection, docs)
		size += len(docs)
		fmt.Printf("Processed %d items\n", size)
		docs = make([]interface{}, 0, blockSize)
	}

	for {
		select {
		case doc := <-c:
			docs = append(docs, doc)

			if len(docs) == blockSize {
				handleDocs()
			}

		case <-quit:
			// Send leftovers
			if len(docs) > 0 {
				handleDocs()
			}
			fmt.Println("Quit")
			return
		}
	}
}

func start(r io.Reader) {
	mongo := NewMongo(*args.mongoURI)
	defer mongo.Close()

	c := make(chan interface{})
	quit := make(chan struct{})

	db := *args.db
	collection := *args.collection
	key := *args.key
	blockSize := *args.blockSize

	fmt.Printf("Db: %s - Collection: %s - Key: %s - BlockSize: %d\n", db, collection, key, blockSize)

	go transmit(r, key, c, quit)
	process(c, quit, mongo, db, collection, key, blockSize)
}

func main() {
	parseArgs()
	start(os.Stdin)
}
