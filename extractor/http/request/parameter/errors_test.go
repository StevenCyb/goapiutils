package parameter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrMissingParameter(t *testing.T) {
	t.Parallel()

	key := "field_name"
	expected := fmt.Sprintf("missing required parameter \"%s\"", key)
	actual := MissingParameterError{Key: key}.Error()
	require.Equal(t, expected, actual)
}

func TestErrMalformedParameter(t *testing.T) {
	t.Parallel()

	key := "param_name"
	expected := fmt.Sprintf("malformed value for parameter \"%s\"", key)
	actual := MalformedParameterError{Key: key}.Error()
	require.Equal(t, expected, actual)
}

func TestTypeMismatchError(t *testing.T) {
	t.Parallel()

	key := "int"
	expected := fmt.Sprintf("parameter type mismatched: expected \"%s\"", key)
	actual := TypeMismatchError{key}.Error()
	require.Equal(t, expected, actual)
}
