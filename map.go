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
	"errors"
)

// OrderedMap is a map that maintains the specific order of elements
// Allows to address elements by key (any comparable) at the same time
// provides deterministic iteration order.
// It is type safe using generics for key and value types.
type OrderedMap[K comparable, T any] struct {
	kvMap  map[K]T
	keys   []K
	keyPos map[K]int
}

// NewOrderedMap creates a new OrderedMap
// K - key type, T - value type
func NewOrderedMap[K comparable, T any]() *OrderedMap[K, T] {
	return &OrderedMap[K, T]{
		kvMap:  make(map[K]T),
		keys:   make([]K, 0),
		keyPos: make(map[K]int),
	}
}

func (m *OrderedMap[K, T]) Init() error {
	if m.kvMap != nil {
		return errors.New("map already initialized")
	}

	m.kvMap = make(map[K]T)
	m.keys = make([]K, 0)
	m.keyPos = make(map[K]int)
	return nil
}

// Set adds a new key-value pair to the map
func (m *OrderedMap[K, T]) Set(key K, value T) {
	if _, ok := m.kvMap[key]; !ok {
		m.keys = append(m.keys, key)    // keep the order, add to the end
		m.keyPos[key] = len(m.keys) - 1 // keep the address of the key
	}
	// Allows rewriting the value
	m.kvMap[key] = value
}

// Get retrieves a value by key
func (m *OrderedMap[K, T]) Get(key K) (any, bool) {
	v, ok := m.kvMap[key]
	return v, ok
}

// Geti retrieves a value by index
func (m *OrderedMap[K, T]) Geti(index int) T {
	return m.kvMap[m.keys[index]]
}

// Getpos retrieves a position by key
func (m *OrderedMap[K, T]) Getpos(key K) int {
	val, ok := m.keyPos[key]
	if !ok {
		return -1
	}
	return val
}

// Keys returns all keys in the map
func (m *OrderedMap[K, T]) Keys() []K {
	return m.keys
}

// Iterator Iterates over the map in the proper order
func (m *OrderedMap[K, T]) Iterator(yield func(key K, value T) bool) {
	for _, key := range m.keys {
		yield(key, m.kvMap[key])
	}
}

// Delete removes a key-value pair from the map
func (m *OrderedMap[K, T]) Delete(key K) {
	delete(m.kvMap, key)

	nilKeyPos := m.keyPos == nil // if keyPos map is nil, we need to recreate it

	_, ok := m.keyPos[key] // check if the key exists
	if ok {
		m.keys = append(m.keys[:m.keyPos[key]], m.keys[m.keyPos[key]+1:]...)
		delete(m.keyPos, key)
	} else { // iterate over the keys to find the key (slow)

		if nilKeyPos {
			m.keyPos = make(map[K]int) // initialize the map
		}

		// since we iterate over all keys we can recreate the keyPos map
		for i, k := range m.keys {
			if k == key {
				m.keys = append(m.keys[:i], m.keys[i+1:]...)
			} else {
				if nilKeyPos {
					m.keyPos[k] = i
				}
			}
		}
	}
}

func (m *OrderedMap[K, T]) Reindex() {
	m.indexKeys()
}

// indexKeys rebuilds the keyPos map
func (m *OrderedMap[K, T]) indexKeys() {
	if m.keyPos == nil {
		m.keyPos = make(map[K]int)
	}
	for i, k := range m.keys {
		m.keyPos[k] = i
	}
}

// Len returns the number of elements in the map
func (m *OrderedMap[K, T]) Len() int {
	return len(m.keys)
}

// Sort by the full map (K-V pair)
func (m *OrderedMap[K, T]) Sort(sortFunc func(map[K]T) ([]K, map[K]int)) {
	m.keys, m.keyPos = sortFunc(m.kvMap)
}

// Sortk sort by keys only
func (m *OrderedMap[K, T]) Sortk(sortFunc func(keys []K) ([]K, map[K]int)) {
	m.keys, m.keyPos = sortFunc(m.keys)
}
