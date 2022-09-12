package patchoperation

import (
	"fmt"
	"reflect"
	"regexp"
)

// Policy specifies the interface for an policy.
type Policy interface {
	GetDetails() string
	Test(operationSpec OperationSpec) bool
}

// DisallowPathPolicy specifies a path that is not allowed.
type DisallowPathPolicy struct {
	Details string
	Path    Path
}

// GetDetails returns the name of this policy.
func (d DisallowPathPolicy) GetDetails() string {
	return d.Details
}

// Test if given operation specification is valid or not.
func (d DisallowPathPolicy) Test(operationSpec OperationSpec) bool {
	return operationSpec.Path != d.Path
}

// DisallowOperationOnPathPolicy disallows specified operation on path.
type DisallowOperationOnPathPolicy struct {
	Details   string
	Path      Path
	Operation Operation
}

// GetDetails returns the name of this policy.
func (d DisallowOperationOnPathPolicy) GetDetails() string {
	return d.Details
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
	Details string
	Path    Path
	Kind    reflect.Kind
}

// GetDetails returns the name of this policy.
func (f ForceTypeOnPathPolicy) GetDetails() string {
	return f.Details
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
	Details    string
	Path       Path
	Expression regexp.Regexp
}

// GetDetails returns the name of this policy.
func (m ForceRegexMatchPolicy) GetDetails() string {
	return m.Details
}

// Test if given operation specification is valid or not.
func (m ForceRegexMatchPolicy) Test(operationSpec OperationSpec) bool {
	if operationSpec.Path != m.Path {
		return true
	}

	return m.Expression.MatchString(fmt.Sprintf("%+v", operationSpec.Value))
}
