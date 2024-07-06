# Ordered Map for Go (Golang)

This library uses an experimental feature for iterators.

You need to enable `GOEXPERIMENT=rangefunc` in order to build it.

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