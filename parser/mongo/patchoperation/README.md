# Query for operational path
With "patch operations" patches can be specified in detail base on [RFC6902](https://datatracker.ietf.org/doc/html/rfc6902)
(`test` operation is not implemented - use [RSQL Parser](parser/mongo/rsql/README.md) instead)
There are five operations available:
| Operation | Description                                                                   |
|-----------|-------------------------------------------------------------------------------|
| `remove`  | remove the value at the target location.                                      |
| `add`     | add a value or array to an array at the target location.                      |
| `replace` | replaces the value at the target location with a new value.                   |
| `move`    | removes the value at a specified location and adds it to the target location. |
| `copy`    | copies the value from a specified location to the target location.            |

This features are supported for `MongoDB 4.2+`.

Additionally, simple rules can be set:
| Policy                          | Description                                                                   |
|---------------------------------|----------------------------------------------------------|
| `DisallowPathPolicy`            | specifies a path that is not allowed.                    |
| `DisallowOperationOnPathPolicy` | disallows specified operation on path.                   |
| `ForceTypeOnPathPolicy`         | forces the value of a specif path to be from given type. |
| `ForceRegexMatchPolicy`         | forces the value of a specif path to match expression.   |

# How to
## Basic usage
```go
import (
  "github.com/StevenCyb/goquery/parser/mongo/patchoperation"

  "go.mongodb.org/mongo-driver/mongo/options"
)
// ...

  var operations []patchoperation.OperationSpec
  err = json.NewDecoder(req.Body).Decode(&operations)
  // ...

  parser := patchoperation.NewParser()
  query, err := parser.Parse(operations...)
  // ...

  result := collection.FindOneAndUpdate(ctx, filter, query, updateOptions)
  // ...
```
## Policy usage
```go
  parser := patchoperation.NewParser(
    DisallowPathPolicy{Details: "illegal ID modification", Path: "_id"},
    ForceTypeOnPathPolicy{Details: "age as number", Path: "user.age", Kind: reflect.Int64},
  ),
```