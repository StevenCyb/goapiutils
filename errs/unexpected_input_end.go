package errs

import (
	"fmt"
)

const errUnexpectedInputEndMessage = "Unexpected end of input, expected: \"%s\""

// ErrUnexpectedInputEnd is an error
// type for unexpected input end
type ErrUnexpectedInputEnd struct {
	tokenType string
}

// Error returns the error message text
func (err ErrUnexpectedInputEnd) Error() string {
	return fmt.Sprintf(errUnexpectedInputEndMessage, err.tokenType)
}

// NewErrUnexpectedInputEnd cerate a new error
func NewErrUnexpectedInputEnd(tokenType string) ErrUnexpectedInputEnd {
	return ErrUnexpectedInputEnd{tokenType: tokenType}
}
