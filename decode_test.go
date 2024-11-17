package omap

import (
	"strings"
	"testing"
)

func TestOrderedMap_DecodeJSON(t *testing.T) {
	data := `{"a": 1, "b": 2, "c": 3, "z": 4,
		"embedded": { "g": 5, "h": 6, "i": 7 }, 
		"arr": [8,9,10]}`

	dataMap := NewOrderedMap[string, any]()

	err := dataMap.Decode(strings.NewReader(data))
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	t.Run("decode JSON", func(t *testing.T) {
		if dataMap.Len() != 6 {
			t.Errorf("Expected 6 elements, got %d", dataMap.Len())
		}

		val, ok := dataMap.Get("z")
		if !ok {
			t.Errorf("Expected key 'z' to be 4, got %v", val)
		}

		if v, ok := val.(float64); !ok || v != 4 {
			t.Errorf("Expected value 4, got %v", val)
		}
	})

	t.Run("check order of elements", func(t *testing.T) {
		i := 1
		for _, v := range dataMap.Iterator {
			switch val := v.(type) {
			case float64:
				vint := int(val)
				if vint != i {
					t.Errorf("Expected value %d, got %v", i, v)
				}
				i++
			case *OrderedMap[string, any]:
				for _, vv := range val.Iterator {
					vint := int(vv.(float64))
					if vint != i {
						t.Errorf("Expected value %d, got %v", i, vv)
					}
					i++
				}
			case []any:
				for _, v := range val {
					vint := int(v.(float64))
					if vint != i {
						t.Errorf("Expected value %d, got %v", i, v)
					}
					i++
				}
			default:
				t.Errorf("Unexpected type %T", val)
			}
		}
	})
}

func TestOrderedMap_DecodeJSONError(t *testing.T) {
	t.Run("invalid JSON", func(t *testing.T) {
		data := `a{"a": 1}`
		dataMap := NewOrderedMap[string, any]()

		err := dataMap.Decode(strings.NewReader(data))
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("missing closing bracket", func(t *testing.T) {
		data := `{"a": 1`
		dataMap := NewOrderedMap[string, any]()

		err := dataMap.Decode(strings.NewReader(data))
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("start", func(t *testing.T) {
		data := `[a]{"a": 1}`
		dataMap := NewOrderedMap[string, any]()

		err := dataMap.Decode(strings.NewReader(data))
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("delimiter", func(t *testing.T) {
		data := `{"a"; 1}`
		dataMap := NewOrderedMap[string, any]()

		err := dataMap.Decode(strings.NewReader(data))
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("key", func(t *testing.T) {
		data := `{a; 1}`
		dataMap := NewOrderedMap[string, any]()

		err := dataMap.Decode(strings.NewReader(data))
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("embedded object", func(t *testing.T) {
		data := `{"a": 1, "obj": {"b"":2}}`
		dataMap := NewOrderedMap[string, any]()

		err := dataMap.Decode(strings.NewReader(data))
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("malformed array", func(t *testing.T) {
		data := `{"a": 1, "ar": ["1"\]}`
		dataMap := NewOrderedMap[string, any]()

		err := dataMap.Decode(strings.NewReader(data))
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("missing array closing bracket", func(t *testing.T) {
		data := `{"a": 1, "ar": ["1"}`
		dataMap := NewOrderedMap[string, any]()

		err := dataMap.Decode(strings.NewReader(data))
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}

func TestOrderedMap_Types(t *testing.T) {
	t.Run("mismatch key type", func(t *testing.T) {
		data := `{"a": 1}`
		dataMap := NewOrderedMap[int, any]()
		err := dataMap.Decode(strings.NewReader(data))
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})

	t.Run("mismatch data type", func(t *testing.T) {
		data := `{"a": 1}`
		dataMap := NewOrderedMap[string, string]()
		err := dataMap.Decode(strings.NewReader(data))
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
