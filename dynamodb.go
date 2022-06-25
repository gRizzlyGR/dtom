package main

import (
	"strconv"
)

type DynamoDBJSONMap map[string]interface{}
type DynamoDBJSONList []interface{}

type DynamoDBItem struct {
	Item *DynamoDBJSONMap
}

func (m *DynamoDBJSONMap) ToValue() JSONValue {
	for key, val := range *m {
		switch key {
		case "S", "BOOL":
			return val
		case "N":
			// TODO: parse to int when possible
			num, err := strconv.ParseFloat(val.(string), 64)
			if err != nil {
				panic(err)
			}
			return num
		case "M":
			jsonMap := DynamoDBJSONMap(val.(map[string]interface{}))
			return jsonMap.ToMap()
		case "L":
			jsonList := DynamoDBJSONList(val.([]interface{}))
			return jsonList.ToList()
		default: // We've found a field name
			field := val.(DynamoDBJSONMap)
			return field.ToValue()
		}
	}

	return nil
}

func (l *DynamoDBJSONList) ToList() JSONList {
	converted := make(JSONList, 0)

	for _, value := range *l {
		jsonMap := DynamoDBJSONMap(value.(map[string]interface{}))
		converted = append(converted, jsonMap.ToValue())
	}

	return converted
}

func (m *DynamoDBJSONMap) ToMap() JSONMap {
	converted := JSONMap{}

	for field, value := range *m {
		jsonMap := DynamoDBJSONMap(value.(map[string]interface{}))
		converted[field] = jsonMap.ToValue()
	}

	return converted
}
