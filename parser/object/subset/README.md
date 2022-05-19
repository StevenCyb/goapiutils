# Query for object subset
Documents can become quite large after a while. 
But if only a small part of the document is needed, it creates unnecessary load.
This parser can be used to generate a subset by means of a query.
This would allow the API to return only the required data. E.g. the roles of a user.

## The language
The syntax of this language is simple: `path.field_name=subset_field_name`.
You can concatenate fields to a subset with the separator `,` e.g. `contact.email=email,contact.phone=phone` to get the a subset like `{email: "___", phone: "___"}`.

## Example
```golang
import (
	"github.com/StevenCyb/goquery/object/subset"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
  subsetExpression := r.URL.Query().Get("subset")

  // ...

  parser := subset.NewParser(nil)
  resultDataSubset, err := parser.Parse(subsetExpression, resultData)

  // ...
  
  json.NewEncoder(w).Encode(resultDataSubset)
}
```