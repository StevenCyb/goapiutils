package errs

import (
	"fmt"
)

const errPolicyViolationMessage = "Policy violation, policy disallow \"%s\""

// ErrPolicyViolation is an error
// type for policy violation
type ErrPolicyViolation struct {
	key string
}

// Error returns the error message text
func (err ErrPolicyViolation) Error() string {
	return fmt.Sprintf(errPolicyViolationMessage, err.key)
}

// NewErrUnexpectedInputEnd cerate a new error
func NewErrPolicyViolation(key string) ErrPolicyViolation {
	return ErrPolicyViolation{key: key}
}
