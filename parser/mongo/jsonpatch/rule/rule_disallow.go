package rule

import (
	"reflect"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
)

// DisallowRule uses the tag `jp_disallow` and
// defines if operations on field are disallowed.
type DisallowRule struct {
	disallow bool
}

// Tag returns tag of the rule.
func (d DisallowRule) Tag() string {
	return "jp_disallow"
}

// UseValue initializes the rule for specified field.
func (d *DisallowRule) UseValue(path operation.Path, _ reflect.Kind, instance interface{}, value string) error {
	disallow, err := getBoolIfNotEmpty(value, string(path), d.Tag())
	if err != nil {
		return err
	}

	d.disallow = *disallow

	return nil
}

// Apply rule on given patch operation specification.
func (d DisallowRule) Apply(operationSpec operation.Spec) error {
	if !d.disallow {
		return nil
	}

	return ErrOperationsNotAllowed
}
