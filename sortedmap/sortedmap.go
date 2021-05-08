package sortedmap

import (
	"sort"
)

type (
	Map struct {
		items map[string]interface{}
		order []string
	}
)

func (m *Map) Add(key string, val interface{}) {
	ok := m.insertMap(key, val)
	if ok {
		ok = m.insertSlice(key)
		if !ok {
			panic("should never occur")
		}
	}
}

func (m *Map) Remove(key string) {
	ok := m.deleteMap(key)
	if ok {
		ok := m.deleteSlice(key)
		if !ok {
			panic("should never occur")
		}
	}
}

func (m *Map) Get(key string) (interface{}, bool) {
	val, ok := m.items[key]
	return val, ok
}

// Iter returns a function that can be used to iterate through the map
// in key-sorted order. The key parameter specifies the start position.
func (m *Map) Iter(key string) func() (string, interface{}, bool) {
	var i int
	if key != "" {
		i = sort.SearchStrings(m.order, key)
	}

	return func() (string, interface{}, bool) {
		if i < len(m.order) {
			k := m.order[i]
			v := m.items[k]
			i++
			return k, v, true
		}
		return "", nil, false
	}
}

// insertMap inserts the key and val into the map and returns
// true if  a new key was inserted
func (m *Map) insertMap(key string, val interface{}) bool {
	if m.items == nil {
		m.items = map[string]interface{}{}
	}
	_, ok := m.items[key]
	m.items[key] = val
	return !ok
}

// insertSlice inserts the key into the slice and returns
// true if a new key was inserted
func (m *Map) insertSlice(key string) bool {
	i := sort.SearchStrings(m.order, key)
	if i == len(m.order) {
		m.order = append(m.order, key)
		return true
	}
	if key == m.order[i] {
		return false
	}
	m.order = append(m.order[:i+1], m.order[i:]...)
	m.order[i] = key
	return true
}

// deleteMap deletes the key from the the map and returns
// true if the key was deleted
func (m *Map) deleteMap(key string) bool {
	_, ok := m.items[key]
	delete(m.items, key)
	return ok
}

// deleteSlice deletes the key from the slice and returns
// true if it was deleted
func (m *Map) deleteSlice(key string) bool {
	i := sort.SearchStrings(m.order, key)
	if i == len(m.order) {
		return false
	}
	if m.order[i] != key {
		return false
	}
	m.order = append(m.order[:i], m.order[i+1:]...)
	return true
}
