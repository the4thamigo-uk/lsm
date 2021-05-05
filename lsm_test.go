package lsm

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Set(t *testing.T) {
	d, err := New(".")
	t.Cleanup(func() {
		require.NoError(t, d.Close(true))
	})
	require.NoError(t, err)
	err = d.Set("key_a", "val_a")
	require.NoError(t, err)
	err = d.Set("key_b", "val_b")
	require.NoError(t, err)
	val, err := d.Get("key_a")
	require.NoError(t, err)
	require.Equal(t, "val_a", val)
	val, err = d.Get("key_b")
	require.NoError(t, err)
	require.Equal(t, "val_b", val)
	val, err = d.Get("key_a")
	require.NoError(t, err)
	require.Equal(t, "val_a", val)
	_, err = d.Get("key_c")
	require.Error(t, err)
}
