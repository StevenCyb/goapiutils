package parameter

import "fmt"

// MissingParameterError is an error type for missing parameter.
type MissingParameterError struct {
	Key string
}

// Error returns the error message text.
func (err MissingParameterError) Error() string {
	return fmt.Sprintf("missing required parameter \"%s\"", err.Key)
}

// MalformedParameterError is an error type for malformed parameter.
type MalformedParameterError struct {
	Key string
}

// Error returns the error message text.
func (err MalformedParameterError) Error() string {
	return fmt.Sprintf("malformed value for parameter \"%s\"", err.Key)
}

// TypeMismatchError is an error for type mismatch.
type TypeMismatchError struct {
	ExpectedType string
}

// Error returns the error message text.
func (err TypeMismatchError) Error() string {
	return fmt.Sprintf("parameter type mismatched: expected \"%s\"", err.ExpectedType)
}
