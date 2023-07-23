package fp

func ReduceIndexed[tA, tB any](f func(tB, tA, int) tB, z tB, a []tA) tB {
	acc := z
	for i, v := range a {
		acc = f(acc, v, i)
	}

	return acc
}

func Reduce[tA, tB any](f func(tB, tA) tB, z tB, a []tA) tB {
	return ReduceIndexed(
		func(acc tB, e tA, _ int) tB {
			return f(acc, e)
		}, z, a)
}

func MapIndexed[tA, tB any](f func(tA, int) tB, a []tA) []tB {
	return ReduceIndexed(
		func(acc []tB, e tA, i int) []tB {
			return append(acc, f(e, i))
		}, make([]tB, 0, len(a)), a)
}

func Map[tA, tB any](f func(tA) tB, args []tA) []tB {
	return MapIndexed(
		func(e tA, _ int) tB {
			return f(e)
		}, args)
}

func FilterIndexed[tA any](p func(tA, int) bool, a []tA) []tA {
	return ReduceIndexed(
		func(acc []tA, e tA, i int) []tA {
			if p(e, i) {
				return append(acc, e)
			} else {
				return acc
			}
		}, []tA{}, a)
}

func Filter[tA any](p func(tA) bool, a []tA) []tA {
	return FilterIndexed(
		func(e tA, _ int) bool {
			return p(e)
		}, a)
}

type Tuple[tA, tB any] struct {
	A tA
	B tB
}

func Zip[tA, tB any](a []tA, b []tB) []Tuple[tA, tB] {
	if la, lb := len(a), len(b); la < lb {
		return ReduceIndexed(
			func(acc []Tuple[tA, tB], e tA, i int) []Tuple[tA, tB] {
				return append(acc, Tuple[tA, tB]{e, b[i]})
			}, make([]Tuple[tA, tB], 0, la), a)
	} else {
		return ReduceIndexed(
			func(acc []Tuple[tA, tB], e tB, i int) []Tuple[tA, tB] {
				return append(acc, Tuple[tA, tB]{a[i], e})
			}, make([]Tuple[tA, tB], 0, lb), b)
	}
}

type Maybe[tA any] struct {
	Empty bool
	Value tA
}

func FindIndexed[tA any](p func(tA, int) bool, a []tA) Maybe[tA] {
	for i, e := range a {
		if p(e, i) {
			return Maybe[tA]{Value: e}
		}
	}
	return Maybe[tA]{Empty: true}
}

func Find[tA any](p func(tA) bool, a []tA) Maybe[tA] {
	return FindIndexed(
		func(e tA, _ int) bool {
			return p(e)
		}, a)
}

func All[tA comparable](p func(tA) bool, a []tA) bool {
	return Find(func(e tA) bool {
		return !p(e)
	}, a).Empty
}

// Basically, this's an alias for golang.org/x/exp/slices Contains(), but
// for some reason no All()-like function was added there so I decided to make
// this slick duplicate for the sake of completeness
func Any[tA comparable](p func(tA) bool, a []tA) bool {
	return !Find(func(e tA) bool {
		return p(e)
	}, a).Empty
}

type RealNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type ComplexNumber interface {
	RealNumber | ~complex64 | ~complex128
}

func Eq[tA comparable](a tA) func(b tA) bool {
	return func(b tA) bool {
		return a == b
	}
}

func NEq[tA comparable](a tA) func(b tA) bool {
	return func(b tA) bool {
		return a != b
	}
}

func Lt[tA RealNumber | string](a tA) func(b tA) bool {
	return func(b tA) bool {
		return a < b
	}
}

func LtEq[tA RealNumber | string](a tA) func(b tA) bool {
	return func(b tA) bool {
		return a <= b
	}
}

func Gt[tA RealNumber | string](a tA) func(b tA) bool {
	return func(b tA) bool {
		return a > b
	}
}

func GtEq[tA RealNumber | string](a tA) func(b tA) bool {
	return func(b tA) bool {
		return a >= b
	}
}

// Function isn't total: it'll panic if applied to an empty slice
func Minimum[tA RealNumber | string](a []tA) tA {
	min := a[0]
	for _, e := range a[1:] {
		if e < min {
			min = e
		}
	}
	return min
}

// Function isn't total: it'll panic if applied to an empty slice
func Maximum[tA RealNumber | string](a []tA) tA {
	max := a[0]
	for _, e := range a[1:] {
		if e > max {
			max = e
		}
	}
	return max
}

// Surprisingly, this function can concatenate a slice of strings as well.
// Returns 0 when a is an empty slice
func Sum[tA ComplexNumber](a []tA) tA {
	var sum tA
	for _, e := range a {
		sum += e
	}
	return sum
}

// Function isn't total: it'll panic if applied to an empty slice
// Returns 1 when a is an empty slice
func Product[tA ComplexNumber](a []tA) tA {
	var prod tA = 1
	for _, e := range a {
		prod *= e
	}
	return prod
}
