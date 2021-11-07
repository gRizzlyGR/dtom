package main

import (
	"bufio"
	"flag"
	"io"
	"os"
)

type Args struct {
	mongoURI   *string
	db         *string
	collection *string
	key        *string
}

var args Args

func parseArgs() {
	args = Args{
		mongoURI:   flag.String("mongoURI", "", "MongoDB URI"),
		db:         flag.String("db", "", "MongoDB database"),
		collection: flag.String("collection", "", "MongoDB collection"),
		key:        flag.String("key", "", "DynamoDB partition key that will be used as _id"),
	}

	flag.Parse()
}

func read(r io.Reader) {
	mongo := NewMongo(*args.mongoURI)
	defer mongo.Close()

	docs := make([]interface{}, 0)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ddbJson := Unmarshal(scanner.Bytes())
		jsonMap := ddbJson.ConvertMap()
		docs = append(docs, jsonMap.ToBSON(*args.key))
	}

	mongo.InsertMany(*args.db, *args.collection, docs)
}

func main() {
	parseArgs()
	read(os.Stdin)
}
