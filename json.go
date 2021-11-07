package main

import (
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

type JSONValue interface{}
type JSONMap map[string]JSONValue
type JSONList []JSONValue

func Unmarshal(raw []byte) DynamoDBJSONMap {
	var item DynamoDBItem

	err := json.Unmarshal(raw, &item)
	if err != nil {
		panic(err)
	}

	m := DynamoDBJSONMap(*item.Item)

	return m
}

func (m *JSONMap) ToBSON(partitionKey string) bson.M {
	var doc bson.M
	if id, ok := (*m)[partitionKey]; ok {
		doc = bson.M{
			"_id": id,
		}

		delete(*m, partitionKey)

		for key, val := range *m {
			doc[key] = val
		}

		return doc
	}

	panic(fmt.Sprintf(`Partion key "%s" not found`, partitionKey))
}
