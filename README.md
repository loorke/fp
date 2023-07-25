# Functional programming library built with Go generics

This package introduce some common functional programming functions as well as their indexed versions.
Implemented functions and types are inspired by Haskell Prelude.

- reduce (foldl)
- map
- filter
- zip
- find
- Curried comparisment functions: Eq, NEq, Lt, LtEq, Gt, GtEq
- Zero value and length predicates: LenS, LenM, EmptyS, NotEmptyS, EmptyM, NotEmptyM, Zero, NotZero
- all
- any
- minimum
- maximum
- sum
- product

The package isn't intended to completely implement Prelude, but rather it's an
useful tool for some casual issues like the following:

```go
type StorageTimeModel struct{
    Time time.Time
}

type TransportTimeModel struct {
    Time int64
}

func RPC_GetTimeseries() []TransportTimeModel {
    var series []StorageTimeModel = storage.GetTimeSeries()

    return fp.Map(func(e StorageTimeModel) TransportTimeModel {
        return TransportTimeModel{
            Time: e.Time.Unix(),
        }
    }, series)
    
}
```