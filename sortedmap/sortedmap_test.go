package sortedmap

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_AddRemove(t *testing.T) {
	var m Map

	m.Add("x", 1)
	invariants(t, m)
	val, ok := m.Get("x")
	require.True(t, ok)
	require.Equal(t, 1, val)

	m.Add("y", 2)
	invariants(t, m)
	val, ok = m.Get("y")
	require.True(t, ok)
	require.Equal(t, 2, val)

	m.Remove("x")
	invariants(t, m)
	val, ok = m.Get("x")
	require.False(t, ok)
	require.Nil(t, val)

	val, ok = m.Get("y")
	invariants(t, m)
	require.True(t, ok)
	require.Equal(t, 2, val)

	m.Remove("y")
	invariants(t, m)
	val, ok = m.Get("x")
	require.False(t, ok)
	require.Nil(t, val)

	val, ok = m.Get("y")
	invariants(t, m)
	require.False(t, ok)
	require.Nil(t, val)

	m.Remove("z")
}

func Test_Iter(t *testing.T) {
	var m Map
	m.Add("z", 3)
	m.Add("w", 1)
	m.Add("y", 2)

	f := func(it func() (string, interface{}, bool)) ([]string, []interface{}) {
		var keys []string
		var vals []interface{}

		for {
			key, val, ok := it()
			if !ok {
				break
			}
			keys = append(keys, key)
			vals = append(vals, val)
		}
		return keys, vals
	}

	keys, vals := f(m.Iter(""))
	invariants(t, m)
	require.Equal(t, []string{"w", "y", "z"}, keys)
	require.Equal(t, []interface{}{1, 2, 3}, vals)

	keys, vals = f(m.Iter("v"))
	invariants(t, m)
	require.Equal(t, []string{"w", "y", "z"}, keys)
	require.Equal(t, []interface{}{1, 2, 3}, vals)

	keys, vals = f(m.Iter("w"))
	invariants(t, m)
	require.Equal(t, []string{"w", "y", "z"}, keys)
	require.Equal(t, []interface{}{1, 2, 3}, vals)

	keys, vals = f(m.Iter("x"))
	invariants(t, m)
	require.Equal(t, []string{"y", "z"}, keys)
	require.Equal(t, []interface{}{2, 3}, vals)

	keys, vals = f(m.Iter("y"))
	invariants(t, m)
	require.Equal(t, []string{"y", "z"}, keys)
	require.Equal(t, []interface{}{2, 3}, vals)

	keys, vals = f(m.Iter("z"))
	invariants(t, m)
	require.Equal(t, []string{"z"}, keys)
	require.Equal(t, []interface{}{3}, vals)

	keys, vals = f(m.Iter("{"))
	invariants(t, m)
	require.Nil(t, keys)
	require.Nil(t, vals)
}

func invariants(t *testing.T, m Map) {
	require.Equal(t, len(m.items), len(m.order))
	for _, key := range m.order {
		_, ok := m.items[key]
		require.True(t, ok)
	}
}
