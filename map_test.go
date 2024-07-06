// MIT License
//
// Copyright (c) 2024 Eli Lap
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package omap

import (
	"testing"
)

func TestOrderedMap_Create(t *testing.T) {
	t.Run("string-int map", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		if m == nil {
			t.Errorf("Expected a map, got nil")
		}
	})

	t.Run("string-string map", func(t *testing.T) {
		m := NewOrderedMap[string, string]()
		if m == nil {
			t.Errorf("Expected a map, got nil")
		}
	})

	t.Run("int-int map", func(t *testing.T) {
		m := NewOrderedMap[int, int]()
		if m == nil {
			t.Errorf("Expected a map, got nil")
		}
	})

	t.Run("int-string map", func(t *testing.T) {
		m := NewOrderedMap[int, string]()
		if m == nil {
			t.Errorf("Expected a map, got nil")
		}
	})

	t.Run("composites map", func(t *testing.T) {
		type KeyStruct struct {
			Name string
			Age  int
		}

		m := NewOrderedMap[KeyStruct, struct{ Age int }]()
		if m == nil {
			t.Errorf("Expected a map, got nil")
		}
	})
}

func TestOrderedMap_Read(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)
	m.Set("d", 4)
	m.Set("e", 5)
	m.Set("f", 6)

	t.Run("retrieve values", func(t *testing.T) {
		if m.Len() != 6 {
			t.Errorf("Expected 6 elements, got %d", m.Len())
		}

		// Retrieve by key
		a, ok := m.Get("a")
		if !ok {
			t.Errorf("Expected key 'a' to be set")
		} else if a != 1 {
			t.Errorf("Expected value 1, got %d", a)
		}

		f, ok := m.Get("f")
		if !ok {
			t.Errorf("Expected key 'f' to be set")
		} else if f != 6 {
			t.Errorf("Expected value 6, got %d", f)
		}

		// Retrieve by index (order)
		if m.Geti(0) != 1 {
			t.Errorf("Expected value 1, got %d", m.Geti(0))
		}

		if val := m.Geti(3); val != 4 {
			t.Errorf("Expected value 4, got %d", val)
		}

		pos := m.Getpos("c")
		if pos != 2 {
			t.Errorf("Expected position 2, got %d", pos)
		}

		posval := m.Geti(pos)
		if posval != 3 {
			t.Errorf("Expected value 3, got %d", posval)
		}
	})

	t.Run("retrieve keys", func(t *testing.T) {
		keys := m.Keys()
		if len(keys) != 6 {
			t.Errorf("Expected 6 keys, got %d", len(keys))
		}

		if keys[0] != "a" {
			t.Errorf("Expected key 'a', got %s", keys[0])
		}

	})
}

func TestOrderedMap_Delete(t *testing.T) {
	t.Run("delete keys", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)
		m.Set("d", 4)
		m.Set("e", 5)
		m.Set("f", 6)

		pos := m.Getpos("c")
		m.Delete("c")
		if m.Len() != 5 {
			t.Errorf("Expected 5 elements, got %d", m.Len())
		}

		if m.Getpos("c") != -1 {
			t.Errorf("Expected key 'c' to be deleted")
		}

		if val := m.Geti(pos); val != 4 {
			t.Errorf("Expected value 4, got %d", val)
		}

	})
}

func TestOrderedMap_Iterator(t *testing.T) {
	m := NewOrderedMap[string, int]()
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)
	m.Set("d", 4)
	m.Set("e", 5)
	m.Set("f", 6)

	t.Run("iterate over keys", func(t *testing.T) {
		var keys []string
		var values []int

		for k, v := range m.Iterator {
			keys = append(keys, k)
			values = append(values, v)
		}

		if len(keys) != 6 {
			t.Errorf("Expected 6 keys, got %d", len(keys))
		}

		if len(values) != 6 {
			t.Errorf("Expected 6 values, got %d", len(values))
		}
	})

	t.Run("stop iteration", func(t *testing.T) {
		m.Iterator(func(key string, value int) bool {
			if key == "c" {
				return false
			}

			return true
		})
	})
}

func TestOrderedMap_Sort(t *testing.T) {
	t.Run("custom sort key function", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)
		m.Set("d", 4)
		m.Set("e", 5)
		m.Set("f", 6)

		// Reverse the order
		m.Sortk(func(keys []string) ([]string, map[string]int) {
			var result []string
			keyPos := map[string]int{}

			for i, k := range keys {
				pos := len(keys) - i
				keyPos[k] = pos
				result = append([]string{k}, result...)
			}

			return result, keyPos
		})

		// Check values were reversed
		var i int
		for _, v := range m.Iterator {
			if v != 6-i {
				t.Errorf("Expected value %d, got %d", 6-i, v)
			}
			i++
		}
	})

	t.Run("custom sort full map function", func(t *testing.T) {
		m := NewOrderedMap[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)
		m.Set("d", 4)
		m.Set("e", 5)
		m.Set("f", 6)

		m.Sort(func(kvMap map[string]int) ([]string, map[string]int) {
			order := []string{"d", "e", "a", "c", "f", "b"}
			keyPos := map[string]int{
				"d": 1,
				"e": 2,
				"a": 3,
				"c": 4,
				"f": 5,
				"b": 6,
			}

			return order, keyPos
		})

		pos := 1
		for k, _ := range m.Iterator {
			if pos == 1 && k != "d" {
				t.Errorf("Expected key 'd', got %s", k)
			}
			if pos == 2 && k != "e" {
				t.Errorf("Expected key 'e', got %s", k)
			}
			if pos == 3 && k != "a" {
				t.Errorf("Expected key 'a', got %s", k)
			}
			if pos == 4 && k != "c" {
				t.Errorf("Expected key 'c', got %s", k)
			}
			if pos == 5 && k != "f" {
				t.Errorf("Expected key 'f', got %s", k)
			}
			if pos == 6 && k != "b" {
				t.Errorf("Expected key 'b', got %s", k)
			}

			pos++
		}

	})
}
