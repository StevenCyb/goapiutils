package patchoperation

import (
	"fmt"
	"reflect"
	"regexp"
)

// Policy specifies the interface for an policy.
type Policy interface {
	GetName() string
	Test(operationSpec OperationSpec) bool
}

// DisallowPathPolicy specifies a path that is not allowed.
type DisallowPathPolicy struct {
	Name string
	Path Path
}

// GetName returns the name of this policy.
func (d DisallowPathPolicy) GetName() string {
	return d.Name
}

// Test if given operation specification is valid or not.
func (d DisallowPathPolicy) Test(operationSpec OperationSpec) bool {
	return operationSpec.Path != d.Path
}

// Test if given operation specification is valid or not.
type DisallowOperationOnPathPolicy struct {
	Name      string
	Path      Path
	Operation Operation
}

// GetName returns the name of this policy.
func (d DisallowOperationOnPathPolicy) GetName() string {
	return d.Name
}

// Test if given operation specification is valid or not.
func (d DisallowOperationOnPathPolicy) Test(operationSpec OperationSpec) bool {
	if operationSpec.Path != d.Path {
		return true
	}

	return d.Operation != operationSpec.Operation
}

// ForceTypeOnPathPolicy forces the value of a specif path to be from given type.
type ForceTypeOnPathPolicy struct {
	Name string
	Path Path
	Kind reflect.Kind
}

// GetName returns the name of this policy.
func (f ForceTypeOnPathPolicy) GetName() string {
	return f.Name
}

// Test if given operation specification is valid or not.
func (f ForceTypeOnPathPolicy) Test(operationSpec OperationSpec) bool {
	if operationSpec.Path != f.Path {
		return true
	}

	return reflect.TypeOf(operationSpec.Value).Kind() == f.Kind
}

// ForceRegexMatchPolicy forces the value of a specif path to match expression.
type ForceRegexMatchPolicy struct {
	Name       string
	Path       Path
	Expression regexp.Regexp
}

// GetName returns the name of this policy.
func (m ForceRegexMatchPolicy) GetName() string {
	return m.Name
}

// Test if given operation specification is valid or not.
func (m ForceRegexMatchPolicy) Test(operationSpec OperationSpec) bool {
	if operationSpec.Path != m.Path {
		return true
	}

	return m.Expression.MatchString(fmt.Sprintf("%+v", operationSpec.Value))
}
