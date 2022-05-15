package errs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrUnexpectedInput(t *testing.T) {
	key := "abc"
	require.Equal(t,
		fmt.Sprintf(errUnexpectedInputMessage, key),
		NewErrUnexpectedInput(key).Error(),
	)
}
