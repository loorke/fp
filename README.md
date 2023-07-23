# Functional programming library built with Go generics

This package introduce some common functional programming functions as well as their indexed versions.
Implemented functions and types mimic Haskell Prelude.

- reduce (foldl)
- map
- filter
- zip
- find
- Curried comparisment functions: Eq, NEq, Lt, LtEq, Gt, GtEq
- all
- any
- minimum
- maximum
- sum
- product

Author didn't intend to write a complete Prelude that could've been used in the same manner as in Haskell.
The package can be useful in some casual situations like the following:

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