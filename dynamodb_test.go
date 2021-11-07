package main

import (
	"reflect"
	"testing"
)

func TestConvertMap(t *testing.T) {
	test := `{"Item":{"key1":{"S":"val1"},"key2":{"N":"42"},"key3":{"BOOL":false},"key4":{"M":{"key4.1":{"S":"val4.1"},"key4.2":{"M":{"key4.2.1":{"S":"val4.2.1"}}}}},"key5":{"L":[{"M":{"key5.0":{"S":"val5.0"},"key5.1":{"N":"42"}}}]},"key6":{"M":{"key6.1":{"L":[{"L":[{"M":{"key6.1.1":{"S":"val6.1.1"}}}]}]}}},"key7":{"L":[]},"key8":{"S":"val8€"}}}`

	m := Unmarshal([]byte(test))
	got := m.ConvertMap()

	want := JSONMap{
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

	if !reflect.DeepEqual(want, got) {
		t.Errorf("ConvertMap(%v):\n\tgot: %v\n\twant: %v", test, got, want)
	}
}
