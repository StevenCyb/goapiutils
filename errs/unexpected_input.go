package errs

import (
	"fmt"
)

const errUnexpectedInputMessage = "Unexpected input \"%v\""

// ErrUnexpectedInput is an error
// type for unexpected input
type ErrUnexpectedInput struct {
	data interface{}
}

// Error returns the error message text
func (err ErrUnexpectedInput) Error() string {
	return fmt.Sprintf(errUnexpectedInputMessage, err.data)
}

// NewErrUnexpectedInput cerate a new error
func NewErrUnexpectedInput(data interface{}) ErrUnexpectedInput {
	return ErrUnexpectedInput{data: data}
}
