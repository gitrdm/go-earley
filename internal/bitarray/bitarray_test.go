package bitarray_test

import (
	"testing"

	"github.com/patrickhuber/go-earley/internal/bitarray"
	"github.com/stretchr/testify/require"
)

func TestBitArray(t *testing.T) {
	ba := bitarray.New(64)
	require.NoError(t, ba.Set(10, true))

	b, err := ba.Get(10)
	require.NoError(t, err)
	require.True(t, b)

	b, err = ba.Get(11)
	require.NoError(t, err)
	require.False(t, b)
}

func TestEdge(t *testing.T) {

	ba := bitarray.New(65)
	require.NoError(t, ba.Set(64, true))

	b, err := ba.Get(64)
	require.NoError(t, err)
	require.True(t, b)

	b, err = ba.Get(63)
	require.NoError(t, err)
	require.False(t, b)
}
