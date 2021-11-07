package main

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestToBSON(t *testing.T) {
	test := JSONMap{
		"key1": "val1",
		"key2": 42.0,
		"key3": false,
		"key4": JSONMap{
			"key4.1": "val4.1",
			"key4.2": JSONMap{
				"key4.2.1": "val4.2.1",
			},
		},
		"key5": JSONList{
			JSONMap{
				"key5.0": "val5.0",
				"key5.1": 42.0,
			},
		},
		"key6": JSONMap{
			"key6.1": JSONList{
				JSONList{
					JSONMap{
						"key6.1.1": "val6.1.1",
					},
				},
			},
		},
		"key7": JSONList{},
		"key8": "val8€",
	}

	got := test.ToBSON("key1")

	want := bson.M{
		"_id":  "val1",
		"key2": 42.0,
		"key3": false,
		"key4": JSONMap{
			"key4.1": "val4.1",
			"key4.2": JSONMap{
				"key4.2.1": "val4.2.1",
			},
		},
		"key5": JSONList{
			JSONMap{
				"key5.0": "val5.0",
				"key5.1": 42.0,
			},
		},
		"key6": JSONMap{
			"key6.1": JSONList{
				JSONList{
					JSONMap{
						"key6.1.1": "val6.1.1",
					},
				},
			},
		},
		"key7": JSONList{},
		"key8": "val8€",
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("ToBSON(%v):\n\tgot: %v\n\twant: %v", test, got, want)
	}
}
