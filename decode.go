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
func (m *OrderedMap[string, any]) Decode(r io.Reader) error {
	dec := json.NewDecoder(r)

	t, err := dec.Token()
	if err != nil {
		return err
	}

	if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expected '{' but got %v", t)
	}

	if err := decodeObject(dec, m); err != nil {
		return err
	}

	return nil
}

func decodeObject[K comparable, T any](dec *json.Decoder, m *OrderedMap[K, T]) error {
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
		if err != nil {
			return err
		}
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

func decodeValue(dec *json.Decoder) (any, error) {
	t, err := dec.Token()
	if err != nil {
		return nil, err
	}

	switch v := t.(type) {
	case json.Delim:
		switch v {
		case '{':
			subMap := NewOrderedMap[string, any]()
			if err := decodeObject(dec, subMap); err != nil {
				return nil, err
			}
			return subMap, nil
		case '[':
			return decodeArray(dec)
		}
	}
	return t, nil
}
