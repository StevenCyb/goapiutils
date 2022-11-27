package rule

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
)

// MinRule uses the tag `jp_min` and defines the minimum size/value:
/*
 * Applies on len for Array, Chan, Map, Slice and String.
 * Applies on value for any numeric type.
 */
type MinRule struct {
	min float64
}

// Tag returns tag of the rule.
func (m MinRule) Tag() string {
	return "jp_min"
}

// UseValue initializes the rule for specified field.
func (m *MinRule) UseValue(path operation.Path, _ reflect.Kind, instance interface{}, value string) error {
	min, err := getFloat64IfNotEmpty(value, string(path), m.Tag())
	if err != nil {
		return err
	}

	m.min = *min

	return nil
}

// Apply rule on given patch operation specification.
func (m MinRule) Apply(operationSpec operation.Spec) error {
	var (
		kind = reflect.TypeOf(operationSpec.Value).Kind()
		size float64
	)

	switch kind { //nolint:exhaustive
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		size = float64(reflect.ValueOf(operationSpec.Value).Len())
	default:
		converted, err := strconv.ParseFloat(fmt.Sprint(operationSpec.Value), 64)
		if err != nil {
			return fmt.Errorf("converting failed: %w", err)
		}

		size = converted
	}

	if m.min > size {
		return LessThenError{ref: m.min, value: size}
	}

	return nil
}
