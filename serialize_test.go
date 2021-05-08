package lsm

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"github.com/the4thamigo-uk/lsm/sortedmap"
	"testing"
)

func Test_ReadWrite(t *testing.T) {
	var m1 sortedmap.Map
	m1.Add("z", 3)
	m1.Add("w", 1)
	m1.Add("y", 2)

	var b bytes.Buffer
	err := WriteMap(&b, m1)
	require.NoError(t, err)

	m2, err := ReadMap(&b)
	require.NoError(t, err)

	require.Equal(t, m1, m2)
}
