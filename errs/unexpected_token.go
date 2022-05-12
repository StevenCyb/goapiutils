package errs

import "fmt"

const errUnexpectedTokenMessage = "Unexpected token: \"%s\" at position \"%d\""

// ErrUnexpectedToken is an error
// type for unexpected token
type ErrUnexpectedToken struct {
	position int
	token    string
}

// Error returns the error message text
func (err ErrUnexpectedToken) Error() string {
	return fmt.Sprintf(errUnexpectedTokenMessage,
		err.token,
		err.position)
}

// NewErrUnexpectedToken cerate a new error
func NewErrUnexpectedToken(position int, token string) ErrUnexpectedToken {
	return ErrUnexpectedToken{
		position: position,
		token:    token,
	}
}
