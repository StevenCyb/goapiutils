package jsonpath

import (
	"reflect"

	"github.com/StevenCyb/goapiutils/parser/errs"
	"go.mongodb.org/mongo-driver/bson"
)

// NewParser creates a new parser.
func NewParser(policies ...Policy) *Parser {
	return &Parser{
		Policies: policies,
	}
}

// Parser that can parse patch operation to generate mongo queries.
type Parser struct {
	Policies []Policy
}

// Parse given operation spec to generate mongo queries if not violating policies.
func (p Parser) Parse(operationSpecs ...OperationSpec) (bson.A, error) {
	for _, policy := range p.Policies {
		for _, operationSpec := range operationSpecs {
			if !operationSpec.Valid() {
				return nil, errs.NewErrUnexpectedInput(operationSpec)
			}

			if !policy.Test(operationSpec) {
				return nil, errs.NewErrPolicyViolation(policy.GetDetails())
			}
		}
	}

	return p.generateMongoQuery(operationSpecs...)
}

// generateMongoQuery generates the mongo query out of operation spec.
func (p Parser) generateMongoQuery(operationSpecs ...OperationSpec) (bson.A, error) {
	var (
		element bson.M
		query   = bson.A{}
	)

	for _, operationSpec := range operationSpecs {
		switch operationSpec.Operation {
		case RemoveOperation:
			element = bson.M{
				"$unset": string(operationSpec.Path),
			}
		case AddOperation:
			if reflect.TypeOf(operationSpec.Value).Kind() != reflect.Slice {
				operationSpec.Value = []interface{}{operationSpec.Value}
			}

			element = bson.M{
				"$set": bson.M{
					string(operationSpec.Path): bson.M{
						"$concatArrays": bson.A{
							"$" + string(operationSpec.Path),
							operationSpec.Value,
						},
					},
				},
			}
		case ReplaceOperation:
			element = bson.M{
				"$set": bson.M{
					string(operationSpec.Path): operationSpec.Value,
				},
			}
		case MoveOperation:
			query = append(query, bson.M{
				"$set": bson.M{
					string(operationSpec.Path): "$" + string(operationSpec.From),
				},
			})
			element = bson.M{
				"$unset": string(operationSpec.From),
			}
		case CopyOperation:
			element = bson.M{
				"$set": bson.M{
					string(operationSpec.Path): "$" + string(operationSpec.From),
				},
			}
		}

		query = append(query, element)
	}

	return query, nil
}
