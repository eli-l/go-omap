# Ordered Map for Go (Golang)

[![codecov](https://codecov.io/gh/eli-l/go-omap/graph/badge.svg?token=JRRGRQYADI)](https://codecov.io/gh/eli-l/go-omap)

This library uses an experimental feature for iterators.

## Iterator functions:
This library depends on iterator function. In Go v1.23 range-over-func a part of the language. 
For Go v1.22 you need to enable `GOEXPERIMENT=rangefunc` in order to build it.

## Ordered map
Provides a way to use Key-Value pais within a map (strictly typed with generic functions).
at the same time it keeps the order of the keys within the internal slice.

## Usage
```go
// Create a new ordered map, string as key and int as value.
// key can be any comparable type, value can be any type.
m := NewOrderedMap[string, int]()

// Set a value referenced by a key
m.Set("a", 1)

// Get a value referenced by a key
// Compatible with the map usage (ok (bool) indicates if key is present in a map)
v, ok := m.Get("a")

// Get by index (order!)
v := m.Geti(3) // get 3rd element

// Get position of a key
p := m.Getpos("c") // get position of key "c"

// Get keys (ordered)
keys := m.Keys() // slice of []<Key type>

// Delete a key
m.Delete("a") // delete key "a", remove from ordered slice, remove from key order

// Iterate over elements in an order
for key, value := range m {
    fmt.Println(key, value)
}

// Get element count
c := m.Len() // get element count
```

### Sorting
This OrderedMap provides a way to sort the map.
There are 2 functions to sort the map:
Sort() and Sortk()

Both functions return ordered slice of keys and a map of keys to their position in the slice.

Second argument may be nil, but this will affect deletion performance.

In case keyPos is set to nil it can be reindex with Reindex() call. It will iterate over the slice 
of the keys to determine position of each key.

```go

// Sort() - recieves unsorted map as argument with original Key-Value pairs

// sortFunc is a closure that takes a map (UNORDERED!) and return a slice of keys (ORDERED)
// and a map of keys to their position in the slice (for fast deletion)
m.Sort(sortFunc func(map[K]T) ([]K, map[K]int))


// Sortk() - sort function get slice of keys only. Slice is already ordered.

m.Sortk(sortFunc func(keys []K) ([]K, map[K]int))

```

### Decode JSON
This library provides a way to decode ordered map from an arbitrary JSON.
It preserves the order of original keys.
Accepts plain values, embedded JSON objects and arrays.

- Embedded object are decoded as OrderedMap[string, any] as well.
- Embedded Arrays are decoded as []any.

```go
data := `{"a": 1, "b": 2, "c": 3, "z": 4,
		"embedded": { "g": 5, "h": 6, "i": 7 }, 
		"arr": [8,9,10]}`

// it must be a map of [string, any], like regular map for parsing JSON
dataMap := NewOrderedMap[string, any]()

err := dataMap.Decode(strings.NewReader(data))
```

See [tests](./decode_test.go) for more examples.

