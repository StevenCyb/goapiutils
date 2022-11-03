# Parameter for HTTP requests
The supported parameter extractors are explained below.

## Query Parameter
The `FromQuery` extractor function provides a simple way to extract `string`, `int`, `float64` and `bool` values from the query (`www....com/item?query=xyz`).
In addition, `Option` provides additional features explained [below](#parameter-option).

```go
import (
  "github.com/StevenCyb/goapiutils/extractor/httprequest/parameter"
)
// ...

stringValue, err := parameter.FromQuery[string](req, parameter.Option{Key: "stringValue"})
intValue, err := parameter.FromQuery[int](req, parameter.Option{Key: "intValue"})
floatValue, err := parameter.FromQuery[float64](req, parameter.Option{Key: "floatValue"})
boolValue, err := parameter.FromQuery[bool](req, parameter.Option{Key: "boolValue"})
```

## Path Parameter
The `FromPath` extractor function provides a simple way to extract `string`, `int`, and `bool` values from the URL path (`www....com/item/{item_id}`).
In addition, `Option` provides additional features explained [below](#parameter-option).

```go
import (
  "github.com/StevenCyb/goapiutils/extractor/httprequest/parameter"
)
// ...

stringValue, err := parameter.FromPath[string](req, parameter.Option{Key: "stringValue"})
intValue, err := parameter.FromPath[int](req, parameter.Option{Key: "intValue"})
boolValue, err := parameter.FromPath[bool](req, parameter.Option{Key: "boolValue"})
```

## Parameter Option
The `Option` argument for the `From*` extractors are similar.
Available options:
* `Key`[string]: name of the parameter
* `Required`[bool]: mark parameter as required 
* `Default`[string]: default value (if not required) 
* `RegexPattern`[string]: regular expression to validate parameter