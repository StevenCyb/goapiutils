package rule

import (
	"reflect"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
)

// AllowedOperationsRule uses the tag `jp_op_allowed` and
// defines witch patch operations are allowed for the field.
type AllowedOperationsRule struct {
	operations []operation.Operation
}

// Tag returns tag of the rule.
func (a AllowedOperationsRule) Tag() string {
	return "jp_op_allowed"
}

// UseValue initializes the rule for specified field.
func (a *AllowedOperationsRule) UseValue(
	path operation.Path, _ reflect.Kind, instance interface{}, value string,
) error {
	operations, err := getOperationsIfNotEmpty(value, string(path), a.Tag())
	if err != nil {
		return err
	}

	a.operations = *operations

	return nil
}

// Apply rule on given patch operation specification.
func (a AllowedOperationsRule) Apply(operationSpec operation.Spec) error {
	for _, operation := range a.operations {
		if operation == operationSpec.Operation {
			return nil
		}
	}

	return OperationNotAllowedError{operation: operationSpec.Operation}
}
