package main

import (
	"strconv"
)

// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Programming.LowLevelAPI.html
type DynamoDbJsonDataDescriptor = string

// const (
// 	S    DynamoDbJsonDataDescriptor = "S"
// 	N    DynamoDbJsonDataDescriptor = "N"
// 	B    DynamoDbJsonDataDescriptor = "B"
// 	BOOL DynamoDbJsonDataDescriptor = "BOOL"
// 	NULL DynamoDbJsonDataDescriptor = "NULL"
// 	M    DynamoDbJsonDataDescriptor = "M"
// 	L    DynamoDbJsonDataDescriptor = "L"
// 	SS   DynamoDbJsonDataDescriptor = "SS"
// 	NS   DynamoDbJsonDataDescriptor = "NS"
// 	BS   DynamoDbJsonDataDescriptor = "BS"
// )

type DynamoDBJSONMap map[string]interface{}
type DynamoDBJSONList []interface{}

type DynamoDBItem struct {
	Item *DynamoDBJSONMap
}

func (m *DynamoDBJSONMap) ConvertValue() JSONValue {
	for key, val := range *m {
		switch key {
		case "S", "BOOL":
			return val
		case "N":
			num, err := strconv.ParseFloat(val.(string), 64)
			if err != nil {
				panic(err)
			}
			return num
		case "M":
			jsonMap := DynamoDBJSONMap(val.(map[string]interface{}))
			return jsonMap.ConvertMap()
		case "L":
			jsonList := DynamoDBJSONList(val.([]interface{}))
			return jsonList.ConvertList()
		default: // We've found a field name
			field := val.(DynamoDBJSONMap)
			return field.ConvertValue()
		}
	}

	return nil
}

func (l *DynamoDBJSONList) ConvertList() JSONList {
	converted := make(JSONList, 0)

	for _, elem := range *l {
		jsonMap := DynamoDBJSONMap(elem.(map[string]interface{}))
		converted = append(converted, jsonMap.ConvertValue())
	}

	return converted
}

func (m *DynamoDBJSONMap) ConvertMap() JSONMap {
	content := JSONMap{}

	for field, data := range *m {
		jsonMap := DynamoDBJSONMap(data.(map[string]interface{}))
		content[field] = jsonMap.ConvertValue()
	}

	return content
}
