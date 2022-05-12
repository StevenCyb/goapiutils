# RSQL for MongoDB-Filter
RSQL is a query language based on [FIQL (Feed Item Query Language)](https://datatracker.ietf.org/doc/html/draft-nottingham-atompub-fiql-00).
It's an URI-friendly syntax and therefore well suited for API's. 

## The language
### Basics
RSQL supports the following value comparison operators for `"string"`, `number` and `bool`:
| Operator | Description       |
|----------|-------------------|
| ==       | equal             |
| !=       | not-equal         |
| =gt=     | greater-than      |
| =ge=     | greater-than-qual |
| =lt=     | less-than         |
| =le=     | less-than-equal   |
| =sw=     | starts with       |
| =wd=     | ends with         |

In addition it supports the following array comparison operators:
| Operator | Description        |
|----------|--------------------|
| =in=           | contains     |
| =out=          | not-contains |

Multiple comparisons can be combined with operators:
| Operator | Description |
|----------|-------------|
| ;        | Logical AND |
| ,        | Logical OR  |

For more advanced queries, `context` may be helpful.
They can be used by round brackets e.g. `(expression;expression),(expression;expression)`.
### RSQL examples
- Log-Level *panic*, *error* or *warning* => `log_level=in=("panic","error","warning")`
- Female with age from 30 or higher => `gender=="female";age=ge=30`
- Binary XOR (one is *1* and one *0*) => `(a==0;b==1),(a==1;b==0)`
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
    // just allow the requester expression with 
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