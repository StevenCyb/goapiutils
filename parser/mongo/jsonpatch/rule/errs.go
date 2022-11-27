package rule

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
)

var ErrOperationsNotAllowed = errors.New("operations are not allowed")

type LessThenError struct {
	ref   float64
	value float64
}

func (l LessThenError) Error() string {
	return fmt.Sprintf("value is less then specified: '%f' < '%f'", l.value, l.ref)
}

type OperationNotAllowedError struct {
	operation operation.Operation
}

func (o OperationNotAllowedError) Error() string {
	return fmt.Sprintf("operation '%s' not allowed", o.operation)
}

type GreaterThenError struct {
	ref   float64
	value float64
}

func (g GreaterThenError) Error() string {
	return fmt.Sprintf("value is greater then specified: '%f' > '%f'", g.value, g.ref)
}

type UnknownFieldError struct {
	name string
}

func (u UnknownFieldError) Error() string {
	return fmt.Sprintf("unknown field '%s'", u.name)
}

type TypeMismatchError struct {
	forKey   bool
	name     string
	expected reflect.Kind
	actual   reflect.Kind
}

func (t TypeMismatchError) Error() string {
	if t.forKey {
		return fmt.Sprintf("'%s' key has invalid kind '%s', must be '%s'", t.name, t.actual.String(), t.expected.String())
	}

	return fmt.Sprintf("'%s' has invalid kind '%s', must be '%s'", t.name, t.actual.String(), t.expected.String())
}

type ExpressionNotMatchError struct {
	expression string
	value      string
}

func (e ExpressionNotMatchError) Error() string {
	return fmt.Sprintf("expression '%s' not match %s", e.expression, e.value)
}
