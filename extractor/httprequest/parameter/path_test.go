package parameter

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestFromPathExtract(t *testing.T) {
	t.Parallel()

	stringValue := "coffee"
	intValue := 1
	boolValue := true

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	req = mux.SetURLVars(req, map[string]string{
		"stringValue": stringValue,
		"intValue":    fmt.Sprint(intValue),
		"boolValue":   fmt.Sprint(boolValue),
	})

	t.Run("StringValue_Success", func(t *testing.T) {
		t.Parallel()

		actual, err := FromPath[string](req, Option{Key: "stringValue"})
		require.NoError(t, err)
		require.Equal(t, stringValue, actual)
	})

	t.Run("IntValue_Success", func(t *testing.T) {
		t.Parallel()

		actual, err := FromPath[int](req, Option{Key: "intValue"})
		require.NoError(t, err)
		require.Equal(t, intValue, actual)
	})

	t.Run("BoolValue_Success", func(t *testing.T) {
		t.Parallel()

		actual, err := FromPath[bool](req, Option{Key: "boolValue"})
		require.NoError(t, err)
		require.Equal(t, boolValue, actual)
	})
}

func TestFromPathRequired(t *testing.T) {
	t.Parallel()

	stringValue := "tea"

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	req = mux.SetURLVars(req, map[string]string{
		"stringValue": stringValue,
	})

	t.Run("Required_Success", func(t *testing.T) {
		t.Parallel()

		actual, err := FromPath[string](req, Option{Key: "stringValue", Required: true})
		require.NoError(t, err)
		require.Equal(t, stringValue, actual)
	})

	t.Run("Required_Fail", func(t *testing.T) {
		t.Parallel()

		_, err := FromPath[string](req, Option{Key: "not_existing", Required: true})
		require.Error(t, err)
		require.True(t, errors.As(err, &MissingParameterError{}))
	})
}

func TestFromPathDefault(t *testing.T) {
	t.Parallel()

	defaultValue := "hello"
	stringValue := "world"

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	req = mux.SetURLVars(req, map[string]string{
		"stringValue": stringValue,
	})

	t.Run("DefaultNotUsed_Success", func(t *testing.T) {
		t.Parallel()

		actual, err := FromPath[string](req, Option{Key: "stringValue", Default: defaultValue})
		require.NoError(t, err)
		require.Equal(t, stringValue, actual)
	})

	t.Run("DefaultUsed_Success", func(t *testing.T) {
		t.Parallel()

		actual, err := FromPath[string](req, Option{Key: "not_existing", Default: defaultValue})
		require.NoError(t, err)
		require.Equal(t, defaultValue, actual)
	})
}

func TestFromPathRegexPattern(t *testing.T) {
	t.Parallel()

	pattern := `^[a-z]+$`
	stringValue := "name"
	stringNumberValue := "123"

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	req = mux.SetURLVars(req, map[string]string{
		"stringValue":       stringValue,
		"stringNumberValue": stringNumberValue,
	})

	t.Run("Required_Success", func(t *testing.T) {
		t.Parallel()

		actual, err := FromPath[string](req, Option{Key: "stringValue", RegexPattern: pattern})
		require.NoError(t, err)
		require.Equal(t, stringValue, actual)
	})

	t.Run("Required_Fail", func(t *testing.T) {
		t.Parallel()

		_, err := FromPath[string](req, Option{Key: "stringNumberValue", RegexPattern: pattern})
		require.Error(t, err)
		require.True(t, errors.As(err, &MalformedParameterError{}))
	})
}
