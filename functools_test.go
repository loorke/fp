package fp

import (
	"errors"
	"log"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReduce(t *testing.T) {
	{
		res := Reduce(func(a, b int) int {
			return a + b
		}, 1, 1, 2, 3, 4)
		require.Equal(t, 11, res)
	}

	{
		res := Reduce(func(a []int, e int) []int {
			return append(a, e)
		}, []int{}, 1, 2, 3, 4)
		require.Equal(t, []int{1, 2, 3, 4}, res)
	}

	{
		res := Reduce(func(a []int, e int) []int {
			return append(a, e)
		}, []int{})
		require.Equal(t, []int{}, res)
	}
}

func TestMap(t *testing.T) {
	{
		res := Map(strconv.Itoa, 1, 2, 3)
		require.Equal(t, []string{"1", "2", "3"}, res)
	}

	{
		res := Map(strconv.Itoa)
		require.Equal(t, []string{}, res)
	}
}

func TestFilter(t *testing.T) {
	{
		res := Filter(func(a int) bool { return a > 2 }, 1, 2, 3, 4, 5)
		require.Equal(t, []int{3, 4, 5}, res)
	}

	{
		res := Filter(func(a int) bool { return a > 2 })
		require.Equal(t, []int{}, res)
	}
}

func TestZip(t *testing.T) {
	{
		res := Zip([]int{1, 2, 3, 4}, []string{"1", "2", "3"})
		require.Equal(t, []Tuple[int, string]{
			{1, "1"},
			{2, "2"},
			{3, "3"},
		}, res)
	}

	{
		res := Zip([]int{1, 2, 3, 4}, []string{})
		require.Equal(t, []Tuple[int, string]{}, res)
	}

	{
		res := Zip([]int{}, []string{"1", "2", "3"})
		require.Equal(t, []Tuple[int, string]{}, res)
	}
}

func TestFind(t *testing.T) {
	{
		res, ok := Find(Eq(3), 1, 2, 3, 4, 5)
		require.True(t, ok)
		require.Equal(t, 3, res)
	}

	{
		res, ok := Find(Eq(666), 1, 2, 3, 4, 5)
		require.False(t, ok)
		require.Zero(t, res)
	}
}

func TestAll(t *testing.T) {
	{
		res := All(Lt(1), 2, 3, 4, 5, 6)
		require.True(t, res)
	}

	{
		res := All(Lt(1), 2, 3, 4, 1, 5, 6)
		require.False(t, res)
	}
}

func TestAny(t *testing.T) {
	{
		res := Any(GtEq(1), 2, 3, 4, 5, 6)
		require.False(t, res)
	}

	{
		res := Any(GtEq(1), 2, 3, 4, 1, 5, 6)
		require.True(t, res)
	}
}

func TestMinMax(t *testing.T) {
	{
		res := Maximum[float64](1, 2, 3, 4, 5)
		require.Equal(t, 5.0, res)
	}

	{
		res := Minimum(1, 2, 3, 4, 5)
		require.Equal(t, 1, res)
	}
}

func TestSumProd(t *testing.T) {
	{
		res := Sum(1, 2, 3, 4, 5)
		require.Equal(t, 15, res)
	}

	{
		res := Sum[int]()
		require.Equal(t, 0, res)
	}

	{
		res := Product[int]()
		require.Equal(t, 1, res)
	}

	{
		res := Product(1, 2, 3, 4, 5)
		require.Equal(t, 120, res)
	}
}

func TestNilAssertions(t *testing.T) {
	func() {
		defer func() {
			v := recover()
			require.NotNil(t, v)
			err, ok := v.(error)
			log.Println(err)
			require.True(t, ok)
			require.True(t, errors.Is(err, ErrNil))
		}()

		var i *int
		MustNonNil(i)
	}()

	func() {
		defer func() {
			v := recover()
			require.NotNil(t, v)
			err, ok := v.(error)
			log.Println(err)
			require.True(t, ok)
			require.True(t, errors.Is(err, ErrNonNil))
		}()

		var i int
		MustNil(&i)
	}()
}
