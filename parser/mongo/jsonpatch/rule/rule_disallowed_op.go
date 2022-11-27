package rule

import (
	"reflect"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
)

// DisallowedOperationsRule uses the tag `jp_op_disallowed` and
// defines witch patch operations are disallowed for the field.
type DisallowedOperationsRule struct {
	operations []operation.Operation
}

// Tag returns tag of the rule.
func (d DisallowedOperationsRule) Tag() string {
	return "jp_op_disallowed"
}

// UseValue initializes the rule for specified field.
func (d *DisallowedOperationsRule) UseValue(
	path operation.Path, _ reflect.Kind, instance interface{}, value string,
) error {
	operations, err := getOperationsIfNotEmpty(value, string(path), d.Tag())
	if err != nil {
		return err
	}

	d.operations = *operations

	return nil
}

// Apply rule on given patch operation specification.
func (d DisallowedOperationsRule) Apply(operationSpec operation.Spec) error {
	for _, operation := range d.operations {
		if operation == operationSpec.Operation {
			return OperationNotAllowedError{operation: operation}
		}
	}

	return nil
}
