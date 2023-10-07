# Go functional programming library built with generics

This library introduces various functional programming miscellaneous that are
useful in a number of scenarious. It aims to both provide some long-awaited 
sugar and examine the limits of recentely introduced Go generics. Implemented 
functions and types are partially inspired by Haskell Prelude.

- reduce (foldl)
- map
- filter
- zip
- find
- Predicates and curried comparisment functions: IsZero, Eq, NEq, Lt, LtEq, Gt, GtEq
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
    }, series...)
    
}
```
