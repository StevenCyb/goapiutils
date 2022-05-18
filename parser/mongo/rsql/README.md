# RSQL for MongoDB-Filter
RSQL is a query language based on [FIQL (Feed Item Query Language)](https://datatracker.ietf.org/doc/html/draft-nottingham-atompub-fiql-00).
It's an URI-friendly syntax and therefore well suited for API's. 

## The language
### Basics
RSQL supports multiple comparison operations.
The following table gives an overview and a matrix that shows which literals can be used with the corresponding operators.
| Operator | Description       | Bool | String | Number | Array | Example                                    |
|----------|-------------------|------|--------|--------|-------|--------------------------------------------|
| ==       | equal             | ✔️   | ✔️    | ✔️     | ❌   | `id=="A43Gd2fd3f32l"`                      |
| !=       | not-equal         | ✔️   | ✔️    | ✔️     | ❌   | `status!="pending"`                        |
| =gt=     | greater-than      | ❌   | ❌    | ✔️     | ❌   | `probability=gt=0.5`                       |
| =ge=     | greater-than-qual | ❌   | ❌    | ✔️     | ❌   | `age=ge=18`                                |
| =lt=     | less-than         | ❌   | ❌    | ✔️     | ❌   | `probability=lt=0.5`                       |
| =le=     | less-than-equal   | ❌   | ❌    | ✔️     | ❌   | `high=le=1.60`                             |
| =sw=     | starts with       | ❌   | ✔️    | ❌     | ❌   | `table=sw="DB_"`                           |
| =ew=     | ends with         | ❌   | ✔️    | ❌     | ❌   | `file=ew=".jpg"`                           |
| =in=     | contains          | ❌   | ❌    | ❌     | ✔️   | `log_level=in=("panic","error","warning")` |
| =out=    | not-contains      | ❌   | ❌    | ❌     | ✔️   | `grade=out=(1,2)`                          |

Multiple comparisons can be combined with composite operators:
| Operator | Description | Example                            |
|----------|-------------|------------------------------------|
| ;        | Logical AND | `gender=="female";age=ge=30`       |
| ,        | Logical OR  | `level=="senior",level=="expert"`  |

For more advanced queries, `context` may be helpful.
They can be used by round brackets e.g. `(expression;expression),(expression;expression)`.
A more accurate example could be a binary XOR (only `a` or `b` is `1`) `(a==0;b==1),(a==1;b==0)`.

## Example
### For API
```golang
import (
	"github.com/StevenCyb/goquery/mongo/rsql"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
  queryExpressionString := r.URL.Query().Get("query")

  parser := rsql.NewParser(nil)
  queryExpression, err := parser.Parse(queryExpressionString)
  // ...

  coll.Find(r.Context(), queryExpression)
  // ...
}
```

### For API with policy
This parser supports two types of policies:
1) `WHITELIST_POLICY` -> disallow everything except given fields
2) `BLACKLIST_POLICY` -> allow everything except given fields
```golang
import (
	"github.com/StevenCyb/goquery/mongo/rsql"
	"github.com/StevenCyb/goquery/tokenizer"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
  queryExpressionString := r.URL.Query().Get("query")

  parser := rsql.NewParser(
    // just allow expression with 
    // "first_name", "last_name" and "age"
    tokenizer.NewPolicy(
      tokenizer.WHITELIST_POLICY,
      "first_name", "last_name", "age",
    )
  )
  queryExpression, err := parser.Parse(queryExpressionString)
  // ...

  coll.Find(r.Context(), queryExpression)
  // ...
}
```
