package omap

import (
	"encoding/json"
	"fmt"
	"io"
)

// Decode a JSON object into an OrderedMap
/*
	The JSON object is expected to be a map of key-value pairs
	This implementation preserves original orders of the keys.
*/
func (m *OrderedMap[K, T]) Decode(r io.Reader) error {
	dec := json.NewDecoder(r)

	t, err := dec.Token()
	if err != nil {
		return err
	}

	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expected '{' but got %v", t)
	}

	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return err
		}
		key, ok := t.(K)
		if !ok {
			return fmt.Errorf("expected string but got %v", t)
		}

		value, err := decodeValue(dec)

		val, ok := value.(T)
		if !ok {
			return fmt.Errorf("value type mismatch: %v", value)
		}
		m.Set(key, val)
	}

	if _, err := dec.Token(); err != nil {
		return err
	}

	return nil
}

func decodeArray(dec *json.Decoder) (any, error) {
	var arr []any

	for dec.More() {
		value, err := decodeValue(dec)
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)
	}

	if _, err := dec.Token(); err != nil {
		return nil, err
	}
	return arr, nil
}

func decodeObject(dec *json.Decoder) (any, error) {
	m := NewOrderedMap[string, any]()

	for dec.More() {
		t, err := dec.Token()
		if err != nil {
			return nil, err
		}

		key, ok := t.(string)
		if !ok {
			return nil, fmt.Errorf("expected string but got %v", t)
		}

		value, err := decodeValue(dec)
		if err != nil {
			return nil, err
		}

		m.Set(key, value)
	}
	if _, err := dec.Token(); err != nil {
		return nil, err
	}

	return m, nil
}

func decodeValue(dec *json.Decoder) (any, error) {
	t, err := dec.Token()
	if err != nil {
		return nil, err
	}

	switch v := t.(type) {
	case json.Delim:
		switch v {
		case '{':
			return decodeObject(dec)
		case '[':
			return decodeArray(dec)
		default:
			return nil, fmt.Errorf("unexpected delimiter %v", v)
		}
	default:
		return t, nil
	}
}
