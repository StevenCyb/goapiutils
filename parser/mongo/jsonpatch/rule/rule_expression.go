package rule

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
)

// ExpressionRule uses the tag `jp_expression` and
// defines expression that needs to be matched by value.
// Anything is printed to string before e.g. int will be a number as string.
type ExpressionRule struct {
	expression string
	regex      regexp.Regexp
}

// Tag returns tag of the rule.
func (e ExpressionRule) Tag() string {
	return "jp_expression"
}

// UseValue initializes the rule for specified field.
func (e *ExpressionRule) UseValue(path operation.Path, _ reflect.Kind, instance interface{}, value string) error {
	regex, err := getRegexpIfNotEmpty(value, string(path), e.Tag())
	if err != nil {
		return err
	}

	e.expression = value
	e.regex = *regex

	return nil
}

// Apply rule on given patch operation specification.
func (e ExpressionRule) Apply(operationSpec operation.Spec) error {
	value := fmt.Sprintf("%+v", operationSpec.Value)

	if !e.regex.Match([]byte(value)) {
		return ExpressionNotMatchError{expression: e.expression, value: value}
	}

	return nil
}
