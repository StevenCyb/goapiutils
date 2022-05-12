package errs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrPolicyViolation(t *testing.T) {
	key := "abc"
	require.Equal(t,
		fmt.Sprintf(errPolicyViolationMessage, key),
		NewErrPolicyViolation(key).Error(),
	)
}
