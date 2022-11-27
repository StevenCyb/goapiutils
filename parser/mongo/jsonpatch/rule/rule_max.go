package rule

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
)

var ErrMaxRuleViolation = errors.New("value greater then specified")

// MaxRule uses the tag `jp_max` and defines the maximum size/value:
/*
 * Applies on len for Array, Chan, Map, Slice and String.
 * Applies on value for any numeric type.
 */
type MaxRule struct {
	max float64
}

// Tag returns tag of the rule.
func (m MaxRule) Tag() string {
	return "jp_max"
}

// UseValue initializes the rule for specified field.
func (m *MaxRule) UseValue(path operation.Path, _ reflect.Kind, instance interface{}, value string) error {
	max, err := getFloat64IfNotEmpty(value, string(path), m.Tag())
	if err != nil {
		return err
	}

	m.max = *max

	return nil
}

// Apply rule on given patch operation specification.
func (m MaxRule) Apply(operationSpec operation.Spec) error {
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

	if m.max < size {
		return GreaterThenError{ref: m.max, value: size}
	}

	return nil
}
