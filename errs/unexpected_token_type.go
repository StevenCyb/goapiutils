package errs

import (
	"fmt"
)

const errUnexpectedTokenTypeMessage = "Unexpected token: \"%s\" at position \"%d\", expected: \"%s\""

// ErrUnexpectedtypes.TokenType is an error
// type for unexpected token type
type ErrUnexpectedTokenType struct {
	position int
	actual   string
	expected string
}

// Error returns the error message text
func (err ErrUnexpectedTokenType) Error() string {
	return fmt.Sprintf(errUnexpectedTokenTypeMessage,
		err.actual, err.position, err.expected)
}

// NewErrUnexpectedTokenType cerate a new error
func NewErrUnexpectedTokenType(position int, actual, expected string) ErrUnexpectedTokenType {
	return ErrUnexpectedTokenType{
		position: position,
		actual:   actual,
		expected: expected,
	}
}
