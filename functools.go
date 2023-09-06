package fp

import (
	"fmt"
)

func ReduceIndexed[tA, tB any](f func(tB, tA, int) tB, z tB, a ...tA) tB {
	acc := z
	for i, v := range a {
		acc = f(acc, v, i)
	}
	return acc
}

func Reduce[tA, tB any](f func(tB, tA) tB, z tB, a ...tA) tB {
	return ReduceIndexed(
		func(acc tB, e tA, _ int) tB {
			return f(acc, e)
		}, z, a...)
}

func ReduceZ[tA, tB any](f func(tB, tA) tB, a ...tA) tB {
	return Reduce(f, Null[tB](), a...)
}

func MapIndexed[tA, tB any](f func(tA, int) tB, a ...tA) []tB {
	return ReduceIndexed(
		func(acc []tB, e tA, i int) []tB {
			return append(acc, f(e, i))
		}, make([]tB, 0, len(a)), a...)
}

func Map[tA, tB any](f func(tA) tB, args ...tA) []tB {
	return MapIndexed(
		func(e tA, _ int) tB {
			return f(e)
		}, args...)
}

func FilterIndexed[tA any](p func(tA, int) bool, a ...tA) []tA {
	return ReduceIndexed(
		func(acc []tA, e tA, i int) []tA {
			if p(e, i) {
				return append(acc, e)
			} else {
				return acc
			}
		}, []tA{}, a...)
}

func Filter[tA any](p func(tA) bool, a ...tA) []tA {
	return FilterIndexed(
		func(e tA, _ int) bool {
			return p(e)
		}, a...)
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
			}, make([]Tuple[tA, tB], 0, la), a...)
	} else {
		return ReduceIndexed(
			func(acc []Tuple[tA, tB], e tB, i int) []Tuple[tA, tB] {
				return append(acc, Tuple[tA, tB]{a[i], e})
			}, make([]Tuple[tA, tB], 0, lb), b...)
	}
}

func FindIndexed[tA any](p func(tA, int) bool, a ...tA) (e tA, ok bool) {
	for i, e := range a {
		if p(e, i) {
			return e, true
		}
	}
	return e, false
}

func Find[tA any](p func(tA) bool, a ...tA) (e tA, ok bool) {
	return FindIndexed(
		func(e tA, _ int) bool {
			return p(e)
		}, a...)
}

func All[tA any](p func(tA) bool, a ...tA) bool {
	_, ok := Find(func(e tA) bool {
		return !p(e)
	}, a...)
	return !ok
}

func Any[tA any](p func(tA) bool, a ...tA) bool {
	_, ok := Find(func(e tA) bool {
		return p(e)
	}, a...)
	return ok
}

// Basically, this's an alias for golang.org/x/exp/slices Contains(), but
// for some reason no All()-like function was added there so I decided to make
// this slick duplicate for the sake of completeness
func Includes[tA comparable](a ...tA) func(tA) bool {
	return func(t tA) bool {
		_, ok := Find(Eq(t), a...)
		return ok
	}
}

type RealNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type ComplexNumber interface {
	RealNumber | ~complex64 | ~complex128
}

type Ordered interface {
	RealNumber | ~string
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

func Lt[tA Ordered](a tA) func(b tA) bool {
	return func(b tA) bool {
		return a < b
	}
}

func LtEq[tA Ordered](a tA) func(b tA) bool {
	return func(b tA) bool {
		return a <= b
	}
}

func Gt[tA Ordered](a tA) func(b tA) bool {
	return func(b tA) bool {
		return a > b
	}
}

func GtEq[tA Ordered](a tA) func(b tA) bool {
	return func(b tA) bool {
		return a >= b
	}
}

func LenS[tA any](l int) func(a []tA) bool {
	return func(a []tA) bool {
		return len(a) == l
	}
}

func EmptyS[tA any](a []tA) bool {
	return len(a) == 0
}

func NotEmptyS[tA any](a []tA) bool {
	return !EmptyS(a)
}

func LenM[tA comparable, tB any](l int) func(a map[tA]tB) bool {
	return func(a map[tA]tB) bool {
		return len(a) == l
	}
}

func EmptyM[tA comparable, tB any](a map[tA]tB) bool {
	return len(a) == 0
}

func NotEmptyM[tA comparable, tB any](a map[tA]tB) bool {
	return !EmptyM(a)
}

func Zero[tA comparable](a tA) bool {
	return a == Null[tA]()
}

func NotZero[tA comparable](a tA) bool {
	return !Zero(a)
}

// Function isn't total: it'll panic if applied to an empty slice
func Minimum[tA Ordered](a ...tA) tA {
	min := a[0]
	for _, e := range a[1:] {
		if e < min {
			min = e
		}
	}
	return min
}

// Function isn't total: it'll panic if applied to an empty slice
func Maximum[tA Ordered](a ...tA) tA {
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
func Sum[tA ComplexNumber](a ...tA) tA {
	var sum tA
	for _, e := range a {
		sum += e
	}
	return sum
}

// Function isn't total: it'll panic if applied to an empty slice
// Returns 1 when a is an empty slice
func Product[tA ComplexNumber](a ...tA) tA {
	var prod tA = 1
	for _, e := range a {
		prod *= e
	}
	return prod
}

func Null[tA any]() tA {
	var z tA
	return z
}

func Assert[tA any](e any) tA {
	return e.(tA)
}

//////////
/// Maps

func ReduceM[tA comparable, tB, tC any](f func(tC, tA, tB) tC, z tC, m map[tA]tB) tC {
	acc := z
	for k, v := range m {
		acc = f(acc, k, v)
	}
	return acc
}

func ReduceMZ[tA comparable, tB, tC any](f func(tC, tA, tB) tC, m map[tA]tB) tC {
	return ReduceM(f, Null[tC](), m)
}

func MapM[tA comparable, tB, tC any](f func(tA, tB) tC, m map[tA]tB) []tC {
	return ReduceM(func(acc []tC, k tA, v tB) []tC {
		return append(acc, f(k, v))
	}, make([]tC, 0, len(m)), m)
}

var ErrNil = errNil("_")

type errNil string

func (e errNil) Error() string {
	return fmt.Sprintf(
		"forbidden nil value of type \"%s\"", string(e))
}

func (e errNil) Is(t error) bool {
	_, ok := t.(errNil)
	return ok
}

func MustNonNil[tA any](v *tA) {
	if v == nil {
		panic(errNil(fmt.Sprintf("%T", v)))
	}
}

func IsNonNil[tA any](v *tA) bool {
	return v != nil
}

var ErrNonNil = errNonNil("_")

type errNonNil string

func (e errNonNil) Error() string {
	return fmt.Sprintf(
		"forbidden non-nil value of type \"%s\"", string(e))
}

func (e errNonNil) Is(t error) bool {
	_, ok := t.(errNonNil)
	return ok
}

func MustNil[tA any](v *tA) {
	if v != nil {
		panic(errNonNil(fmt.Sprintf("%T", v)))
	}
}

func IsNil[tA any](v *tA) bool {
	return v == nil
}
