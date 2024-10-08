/*
This library introduces various functional programming miscellaneous that are
useful in a number of scenarious.

# Examples

Processing collections

	var s []string
	s = fp.Map(strconv.Itoa, 1, 2, 3, 4, 5)

	fp.Filter(fp.Gt(100), 1, 2, 555, 0)

Checks and validations

	arr := []int{1, 2, 2, 2, 5, 11}
	if one2ten := fp.All(
		fp.Includes(1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
		arr...,
	); one2ten {

	}

	var FreeFood []food.Food
	FreeFood = Set(food.Butter, food.Bread, food.Butter) // panics
	FreeFood = Set(food.Butter, food.Bread, food.Oil) // ok
*/

package fp

import (
	"fmt"
	"math/rand"
)

func ReduceIndex[
	tA, tB any,
	tF ~func(tB, tA, int) tB,
](f tF, z tB, a ...tA) tB {
	acc := z
	for i, v := range a {
		acc = f(acc, v, i)
	}
	return acc
}

func Reduce[
	tA, tB any,
	tF ~func(tB, tA) tB,
](f tF, z tB, a ...tA) tB {
	return ReduceIndex(
		func(acc tB, e tA, _ int) tB {
			return f(acc, e)
		}, z, a...)
}

func ReduceZ[
	tA, tB any,
	tF ~func(tB, tA) tB,
](f tF, a ...tA) tB {
	return Reduce(f, Zero[tB](), a...)
}

func MapIndex[
	tA, tB any,
	tF ~func(tA, int) tB,
](f tF, a ...tA) []tB {
	return ReduceIndex(
		func(acc []tB, e tA, i int) []tB {
			return append(acc, f(e, i))
		}, make([]tB, 0, len(a)), a...)
}

func Map[
	tA, tB any,
	tF ~func(tA) tB,
](f tF, args ...tA) []tB {
	return MapIndex(
		func(e tA, _ int) tB {
			return f(e)
		}, args...)
}

func FilterIndex[
	tA any,
	tF ~func(tA, int) bool,
](p tF, a ...tA) []tA {
	return ReduceIndex(
		func(acc []tA, e tA, i int) []tA {
			if p(e, i) {
				return append(acc, e)
			} else {
				return acc
			}
		}, []tA{}, a...)
}

func Filter[
	tA any,
	tF ~func(tA) bool,
](p tF, a ...tA) []tA {
	return FilterIndex(
		func(e tA, _ int) bool {
			return p(e)
		}, a...)
}

func Count[
	tA any,
	tF ~func(tA) bool,
](p tF, a ...tA) int {
	var i int
	for _, e := range a {
		if p(e) {
			i++
		}
	}
	return i
}

type Tuple[tA, tB any] struct {
	A tA
	B tB
}

func Zip[tA, tB any](a []tA, b []tB) []Tuple[tA, tB] {
	if la, lb := len(a), len(b); la < lb {
		return ReduceIndex(
			func(acc []Tuple[tA, tB], e tA, i int) []Tuple[tA, tB] {
				return append(acc, Tuple[tA, tB]{e, b[i]})
			}, make([]Tuple[tA, tB], 0, la), a...)
	} else {
		return ReduceIndex(
			func(acc []Tuple[tA, tB], e tB, i int) []Tuple[tA, tB] {
				return append(acc, Tuple[tA, tB]{a[i], e})
			}, make([]Tuple[tA, tB], 0, lb), b...)
	}
}

func FindIndex[
	tA any,
	tF ~func(tA, int) bool,
](p tF, a ...tA) (e tA, ok bool) {
	for i, e := range a {
		if p(e, i) {
			return e, true
		}
	}
	return e, false
}

func Find[
	tA any,
	tF ~func(tA) bool,
](p tF, a ...tA) (e tA, ok bool) {
	return FindIndex(
		func(e tA, _ int) bool {
			return p(e)
		}, a...)
}

func All[
	tA any,
	tF ~func(tA) bool,
](p tF, a ...tA) bool {
	_, ok := Find(func(e tA) bool {
		return !p(e)
	}, a...)
	return !ok
}

func Any[
	tA any,
	tF ~func(tA) bool,
](p tF, a ...tA) bool {
	_, ok := Find(func(e tA) bool {
		return p(e)
	}, a...)
	return ok
}

func Includes[tA comparable](a ...tA) func(tA) bool {
	return func(t tA) bool {
		_, ok := Find(Eq(t), a...)
		return ok
	}
}

//////////
/// Parameters spreading/shrinking

func Shrink[
	tA, tB any,
	tF ~func(...tA) tB,
](f tF) func([]tA) tB {
	return func(a []tA) tB {
		return f(a...)
	}
}

func Spread[
	tA, tB any,
	tC ~[]tA,
	tF ~func(tC) tB,
](f tF) func(...tA) tB {
	return func(a ...tA) tB {
		return f(a)
	}
}

//////////
/// Ordering and predicates

type IntegerNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type FloatNumber interface {
	~float32 | ~float64
}

type RealNumber interface {
	IntegerNumber | FloatNumber
}

type ComplexNumber interface {
	~complex64 | ~complex128
}

type Number interface {
	RealNumber | ComplexNumber
}

// Resemble cmp.Ordered, but better due to ~rune and ~byte being added.
type Ordered interface {
	RealNumber | ~string | ~rune | ~byte
}

type ByteString interface {
	~string | ~[]byte
}

// Not(p)(x) == !p(x)
func Not[
	tA any,
	tF ~func(tA) bool,
](p tF) tF {
	return func(v tA) bool {
		return !p(v)
	}
}

// b == a
func Eq[tA comparable](a tA) func(b tA) bool {
	return func(b tA) bool {
		return a == b
	}
}

// b != a
func NEq[tA comparable](a tA) func(b tA) bool {
	return func(b tA) bool {
		return b != a
	}
}

// b < a
func Lt[tA Ordered](a tA) func(b tA) bool {
	return func(b tA) bool {
		return b < a
	}
}

// b <= a
func LtEq[tA Ordered](a tA) func(b tA) bool {
	return func(b tA) bool {
		return b <= a
	}
}

// b > a
func Gt[tA Ordered](a tA) func(b tA) bool {
	return func(b tA) bool {
		return b > a
	}
}

// b >= a
func GtEq[tA Ordered](a tA) func(b tA) bool {
	return func(b tA) bool {
		return b >= a
	}
}

func Len[tB any](a ...tB) int {
	return len(a)
}

func Len_[tA any, tS ~[]tA](a tS) int {
	return len(a)
}

func MLen[
	tA RealNumber,
	tM ~map[tB]tC,
	tB comparable,
	tC any,
](m tM) tA {
	return tA(len(m))
}

func IsEmpty[
	tS ~[]tA,
	tA any,
](a tS) bool {
	return len(a) == 0
}

func MIsEmpty[
	tM map[tA]tB,
	tA comparable,
	tB any,
](m tM) bool {
	return len(m) == 0
}

func IsZero[tA comparable](a tA) bool {
	return a == Zero[tA]()
}

func IsNotZero[tA comparable](a tA) bool {
	return a != Zero[tA]()
}

// Returns zero value and -1 if no arguments are provided
func MinimumIndex[tA Ordered](a ...tA) (min tA, indx int) {
	if len(a) == 0 {
		return Zero[tA](), -1
	}

	min = a[0]
	for i, e := range a[1:] {
		if e < min {
			indx = i
			min = e
		}
	}
	return min, indx
}

// Returns zero value if no arguments are provided
func Minimum[tA Ordered](a ...tA) tA {
	min, _ := MinimumIndex(a...)
	return min
}

// Returns zero value and -1 if no arguments are provided
func MaximumIndex[tA Ordered](a ...tA) (max tA, indx int) {
	if len(a) == 0 {
		return Zero[tA](), -1
	}

	max = a[0]
	for i, e := range a[1:] {
		if e > max {
			indx = i
			max = e
		}
	}
	return max, indx
}

// Returns zero value if no arguments are provided
func Maximum[tA Ordered](a ...tA) tA {
	max, _ := MaximumIndex(a...)
	return max
}

// Returns duplicate value dup and its index if there's any,
// zero value and index -1 otherwise
// NOTE: More efficient implementations are possible if tA is Ordered or
// if you're willing to modify original slice.
func FindDupsIndex[tA comparable](a ...tA) (dup tA, index int) {
	d := map[tA]struct{}{}
	for i, v := range a {
		if _, ok := d[v]; ok {
			return v, i
		}
		d[v] = struct{}{}
	}

	return Zero[tA](), -1
}

// Finds duplicates, returns zero value and false if none are present;
// see FindDupsIndex()
func FindDups[tA comparable](a ...tA) (dup tA, ok bool) {
	dup, indx := FindDupsIndex(a...)
	return dup, indx != -1
}

// Returns true if no duplicates are present
func NoDups[tA comparable](a ...tA) bool {
	_, ok := FindDups(a...)
	return !ok
}

//////////
/// Arithmetics

// Adds to numbers
func Add[tA Number | string](a tA) func(tA) tA {
	return Lazy2(Apply2(Sum[tA]))(a)
}

func Mul[tA Number](a tA) func(tA) tA {
	return Apply1(Product[tA])
}

// Adds numbers or concatenates strings.
// Returns zero value if a is an empty slice.
func Sum[tA Number | string](a ...tA) tA {
	var sum tA
	for _, e := range a {
		sum += e
	}
	return sum
}

// Multiplies numbers.
// Returns 1 if a is an empty slice.
func Product[tA Number](a ...tA) tA {
	var prod tA = 1
	for _, e := range a {
		prod *= e
	}
	return prod
}

func IsOdd[tA IntegerNumber](n tA) bool {
	return n%2 != 0
}

func IsEven[tA IntegerNumber](n tA) bool {
	return n%2 == 0
}

//////////
/// Maps

func MReduce[
	tF ~func(tC, tA, tB) tC,
	tM ~map[tA]tB,
	tA comparable,
	tB, tC any,
](f tF, z tC, m tM) tC {

	acc := z
	for k, v := range m {
		acc = f(acc, k, v)
	}
	return acc
}

func MReduceZ[
	tF ~func(tC, tA, tB) tC,
	tM ~map[tA]tB,
	tA comparable,
	tB, tC any,
](
	f tF, m tM) tC {
	return MReduce(f, Zero[tC](), m)
}

func MMapK[
	tF ~func(tA, tB) tC,
	tM ~map[tA]tB,
	tA comparable,
	tB, tC any,
](f tF, m tM) []tC {
	return MReduce(func(acc []tC, k tA, v tB) []tC {
		return append(acc, f(k, v))
	}, make([]tC, 0, len(m)), m)
}

func MMap[
	tF ~func(tB) tC,
	tM ~map[tA]tB,
	tA comparable, tB, tC any,
](f tF, m tM) []tC {
	return MMapK(func(_ tA, v tB) tC {
		return f(v)
	}, m)
}

func MMapMK[
	tF ~func(tA, tC) (tB, tD),
	tM ~map[tA]tC,
	tA, tB comparable,
	tC, tD any,
](f tF, m tM) map[tB]tD {

	return MReduce(func(acc map[tB]tD, k tA, v tC) map[tB]tD {
		nk, nv := f(k, v)
		acc[nk] = nv
		return acc
	}, make(map[tB]tD, len(m)), m)
}

func MMapM[
	tF ~func(tB) tC,
	tM ~map[tA]tB,
	tA comparable,
	tB, tC any,
](f tF, m tM) map[tA]tC {
	return MMapMK(func(k tA, v tB) (tA, tC) {
		return k, f(v)
	}, m)
}

//////////
/// Checks and validations

// Ensures uniqueness of all passed arguments
func Enum[tA comparable](a ...tA) []tA {
	dup, i := FindDupsIndex(a...)
	Must(Eq(-1), fmt.Sprintf("duplicate value \"%v\"; index: %d", dup, i), i)
	return a
}

func Must[
	tF ~func(v tA) bool,
	tA any,
](p tF, exceptionMsg string, a ...tA) {
	for i, e := range a {
		if !p(e) {
			panic(MustError{
				Msg: exceptionMsg,
				Val: e,
				Pos: CondZ(&i)(len(a) > 1),
			})
		}
	}
}

func MustNonNil(a ...any) {
	Must(IsNotZero, "nil value is not allowed", a...)
}

func MustNil(a ...any) {
	Must(IsZero, "nil value is mandatory", a...)
}

func NoError[tA any](v tA, err error) tA {
	MustNil(err)
	return v
}

func NoError2[tA, tB any](a tA, b tB, err error) (tA, tB) {
	MustNil(err)
	return a, b
}

func NoError3[tA, tB, tC any](a tA, b tB, c tC, err error) (tA, tB, tC) {
	MustNil(err)
	return a, b, c
}

type MustError struct {
	Msg string
	Val any
	Pos *int // non-nil if multiple values were provided to Must()
}

func (e MustError) Error() string {
	var pos string
	if e.Pos != nil {
		pos = fmt.Sprintf(" at position %d", *e.Pos)
	}

	return fmt.Sprintf("failure for value \"%v\"%s: %s", e.Val, pos, e.Msg)
}

func (e MustError) Unwrap() error {
	err, ok := e.Val.(error)
	if ok {
		return err
	}
	return nil
}

//////////
/// Conditions

// Returns right if ok is true; left otherwise
func Cond[tA any](left, right tA) func(ok bool) tA {
	return func(ok bool) tA {
		if ok {
			return right
		} else {
			return left
		}
	}
}

// // Returns right if ok is true; zero value otherwise
func CondZ[tA any](right tA) func(ok bool) tA {
	return Cond(Zero[tA](), right)
}

//////////
/// Types

func AssertAll[tA, tB any](a ...tB) []tA {
	return Map(func(v tB) tA {
		return Assert[tA](v)
	}, a...)
}

func Assert[tA any](src any) tA {
	dst := src.(tA)
	return dst
}

func Ref[tA any](v tA) *tA {
	return &v
}

func Zero[tA any]() tA {
	var z tA
	return z
}

func ZeroOf[tA any](tA) tA {
	return Zero[tA]()
}

func CastNum[tA, tB RealNumber](a tB) tA {
	return tA(a)
}

func CastStr[tA, tB ByteString](a tB) tA {
	return tA(a)
}

//////////
/// Random

// Returns zero value and -1 if len(a) == 0
func RandChooseIndex[tA any](a ...tA) (tA, int) {
	if len(a) == 0 {
		return Zero[tA](), -1
	}

	i := rand.Intn(len(a))
	return a[i], i
}

func RandChoose[tA any](a ...tA) tA {
	v, _ := RandChooseIndex(a...)
	return v
}

func RandShuffle[tA any](a ...tA) []tA {
	res := make([]tA, len(a))
	copy(res, a)
	RandShuffle_(res)
	return res
}

func RandShuffle_[tS ~[]tA, tA any](a tS) {
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
}

//////////
/// Concatenations

func Concat[tA any](a ...[]tA) []tA {
	var cap int
	for _, e := range a {
		cap += len(e)
	}

	return Reduce(func(acc, a []tA) []tA {
		return append(acc, a...)
	}, make([]tA, 0, cap), a...)
}

//////////
/// Applicators

func Apply1[tA, tB any, tF ~func(...tA) tB](f tF) func(a tA) tB {
	return func(a tA) tB {
		return f(a)
	}
}

func Apply2[tA, tB any, tF ~func(...tA) tB](f tF) func(a, b tA) tB {
	return func(a, b tA) tB {
		return f(a, b)
	}
}

func Apply3[
	tA, tB any,
	tF ~func(...tA) tB,
](f tF) func(a, b, c tA) tB {
	return func(a, b, c tA) tB {
		return f(a, b, c)
	}
}

func Apply4[
	tA, tB any,
	tF ~func(...tA) tB,
](f tF) func(a, b, c, d tA) tB {
	return func(a, b, c, d tA) tB {
		return f(a, b, c, d)
	}
}

func Apply5[
	tA, tB any,
	tF ~func(...tA) tB,
](f tF) func(a, b, c, d, e tA) tB {
	return func(a, b, c, d, e tA) tB {
		return f(a, b, c, d, e)
	}
}

func Lazy2[
	tA, tB, tC any,
	tF ~func(tA, tB) tC,
](f tF) func(a tA) func(tB) tC {
	return func(a tA) func(tB) tC {
		return func(b tB) tC {
			return f(a, b)
		}
	}
}

func Lazy3[
	tA, tB, tC, tD any,
	tF ~func(tA, tB, tC) tD,
](f tF) func(a tA) func(tB) func(tC) tD {
	return func(a tA) func(tB) func(tC) tD {
		return func(b tB) func(tC) tD {
			return func(c tC) tD {
				return f(a, b, c)
			}
		}
	}
}

func Lazy4[
	tA, tB, tC, tD, tE any,
	tF ~func(tA, tB, tC, tD) tE,
](f tF) func(a tA) func(tB) func(tC) func(tD) tE {
	return func(a tA) func(tB) func(tC) func(tD) tE {
		return func(b tB) func(tC) func(tD) tE {
			return func(c tC) func(tD) tE {
				return func(d tD) tE {
					return f(a, b, c, d)
				}
			}
		}
	}
}

func Lazy5[
	tA, tB, tC, tD, tE, tF any,
	tG ~func(tA, tB, tC, tD, tE) tF,
](f tG) func(a tA) func(tB) func(tC) func(tD) func(tE) tF {
	return func(a tA) func(tB) func(tC) func(tD) func(tE) tF {
		return func(b tB) func(tC) func(tD) func(tE) tF {
			return func(c tC) func(tD) func(tE) tF {
				return func(d tD) func(tE) tF {
					return func(e tE) tF {
						return f(a, b, c, d, e)
					}
				}
			}
		}
	}
}
