package set

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSets(t *testing.T) {
	a := New(1, 2, 3, 4, 5)
	require.False(t, a[6])
	require.True(t, a[5])
	b := New(5, 6, 7, 8, 9, 10)

	{
		u := a.Union(b)
		require.True(t, u[4])
		require.True(t, u[5])
		require.True(t, u[6])
	}

	{
		i := a.Intersection(b)
		require.False(t, i[4])
		require.True(t, i[5])
		require.False(t, i[6])
	}

	{
		d := a.Diff(b)
		require.True(t, d[4])
		require.False(t, d[5])
		require.False(t, d[6])
	}

	{
		d := a.SymmetricDiff(b)
		require.True(t, d[4])
		require.False(t, d[5])
		require.True(t, d[6])
	}
}
