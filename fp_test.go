package fp

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShrinkSpread(t *testing.T) {
	require.Equal(t, 6,
		Spread(Shrink(Product[int]))(1, 2, 3))
}

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

func TestCount(t *testing.T) {
	c := Count(IsNotZero, 1, 2, 3, 4, 5)
	require.Equal(t, 5, c)
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

func TestPredicatesAndOrdering(t *testing.T) {
	type predicate func(int) bool
	var p predicate = func(x int) bool {
		return x == 666
	}

	require.True(t, Not(p)(600))
	require.True(t, All(Gt(10), 11, 22, 30))
	require.True(t, Any(IsZero[int], 12, 13, 0, 14))
}

func TestConditions(t *testing.T) {
	require.True(t, Cond(true, false)(false))
	require.True(t, CondZ(true)(true))
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
		res := All(Gt(1), 2, 3, 4, 5, 6)
		require.True(t, res)
	}

	{
		res := All(Gt(1), 2, 3, 4, 1, 5, 6)
		require.False(t, res)
	}
}

func TestAny(t *testing.T) {
	{
		res := Any(LtEq(1), 2, 3, 4, 5, 6)
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

func TestNoDuplicates(t *testing.T) {
	dup, index := FindDupsIndex(1, 2, 3, 4, 5)
	require.Equal(t, -1, index)
	require.Zero(t, dup)

	dup, index = FindDupsIndex(1, 2, 1, 3, 4, 5)
	require.Equal(t, 2, index)
	require.Equal(t, 1, dup)

	require.False(t, NoDups(1, 2, 1, 3, 4, 5))
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

func TestMust(t *testing.T) {
	func() {
		defer func() {
			if rec := recover(); rec != nil {
				err := rec.(error)
				require.True(t, errors.As(err, &MustError{}))
				me := err.(MustError)
				require.Equal(t,
					"failure for value \"1\": greater than zero",
					me.Error())
			}
		}()

		Must(Lt(0), "greater than zero", 1)
	}()

	func() {
		defer func() {
			if rec := recover(); rec != nil {
				err := rec.(error)
				require.True(t, errors.As(err, &MustError{}))
				me := err.(MustError)
				require.Equal(t,
					"failure for value \"1\" at position 1: greater than zero",
					me.Error())
			}
		}()

		Must(Lt(0), "greater than zero", -1, 1)

		Must(IsNotZero, "ololo", 1, 2, 3)
	}()

	func() {
		defer func() {
			if rec := recover(); rec != nil {
				err := rec.(error)
				require.True(t, errors.As(err, &MustError{}))
				me := err.(MustError)
				require.Equal(t,
					"failure for value \"6\": duplicate value \"6\"; index: 6",
					me.Error())
			}
		}()

		Enum(1, 2, 3, 4, 5, 6, 6)
	}()

	func() {
		err1 := errors.New("wrapped error")
		defer func() {
			if rec := recover(); rec != nil {
				err := rec.(error)
				require.True(t, errors.As(err, &MustError{}))
				require.True(t, errors.Is(err, err1))
			}
		}()

		func() {
			defer func() {
				if rec := recover(); rec != nil {
					err := rec.(error)
					panic(errors.Join(err))
				}
			}()

			Must(IsZero, "Unwrap()", err1)
		}()
	}()

	func() {
		defer func() {
			rec := recover()
			require.Nil(t, rec)
		}()

		MustNonNil("must be not nil", 0)
	}()
}

func TestSet(t *testing.T) {
	require.True(t, Set(1, 2, 3, 4)[3])
}

func TestRef(t *testing.T) {
	f := func(x *string) string {
		return *x
	}

	require.Equal(t, "ololo", f(Ref("ololo")))
}
