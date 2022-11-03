package parameter

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromQueryExtract(t *testing.T) {
	t.Parallel()

	stringValue := "coffee"
	intValue := 1
	floatValue := 2.2
	boolValue := true

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	values := req.URL.Query()
	values.Add("stringValue", stringValue)
	values.Add("intValue", fmt.Sprint(intValue))
	values.Add("floatValue", fmt.Sprint(floatValue))
	values.Add("boolValue", fmt.Sprint(boolValue))
	req.URL.RawQuery = values.Encode()

	t.Run("StringValue_Success", func(t *testing.T) {
		t.Parallel()

		actual, err := FromQuery[string](req, Option{Key: "stringValue"})
		require.NoError(t, err)
		require.Equal(t, stringValue, actual)
	})

	t.Run("IntValue_Success", func(t *testing.T) {
		t.Parallel()

		actual, err := FromQuery[int](req, Option{Key: "intValue"})
		require.NoError(t, err)
		require.Equal(t, intValue, actual)
	})

	t.Run("Float64Value_Success", func(t *testing.T) {
		t.Parallel()

		actual, err := FromQuery[float64](req, Option{Key: "floatValue"})
		require.NoError(t, err)
		require.Equal(t, floatValue, actual)
	})

	t.Run("BoolValue_Success", func(t *testing.T) {
		t.Parallel()

		actual, err := FromQuery[bool](req, Option{Key: "boolValue"})
		require.NoError(t, err)
		require.Equal(t, boolValue, actual)
	})
}

func TestFromQueryRequired(t *testing.T) {
	t.Parallel()

	stringValue := "chicken"

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	values := req.URL.Query()
	values.Add("stringValue", stringValue)
	req.URL.RawQuery = values.Encode()

	t.Run("Required_Success", func(t *testing.T) {
		t.Parallel()

		actual, err := FromQuery[string](req, Option{Key: "stringValue", Required: true})
		require.NoError(t, err)
		require.Equal(t, stringValue, actual)
	})

	t.Run("Required_Fail", func(t *testing.T) {
		t.Parallel()

		_, err := FromQuery[string](req, Option{Key: "not_existing", Required: true})
		require.Error(t, err)
		require.True(t, errors.As(err, &MissingParameterError{}))
	})
}

func TestFromQueryDefault(t *testing.T) {
	t.Parallel()

	defaultValue := "hi"
	stringValue := "there"

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	values := req.URL.Query()
	values.Add("stringValue", stringValue)
	req.URL.RawQuery = values.Encode()

	t.Run("DefaultNotUsed_Success", func(t *testing.T) {
		t.Parallel()

		actual, err := FromQuery[string](req, Option{Key: "stringValue", Default: defaultValue})
		require.NoError(t, err)
		require.Equal(t, stringValue, actual)
	})

	t.Run("DefaultUsed_Success", func(t *testing.T) {
		t.Parallel()

		actual, err := FromQuery[string](req, Option{Key: "not_existing", Default: defaultValue})
		require.NoError(t, err)
		require.Equal(t, defaultValue, actual)
	})
}

func TestFromQueryRegexPattern(t *testing.T) {
	t.Parallel()

	pattern := `^[a-z]+$`
	stringValue := "soup"
	stringNumberValue := "123"

	req, err := http.NewRequest("GET", "/", nil)
	require.NoError(t, err)

	values := req.URL.Query()
	values.Add("stringValue", stringValue)
	values.Add("stringNumberValue", stringNumberValue)
	req.URL.RawQuery = values.Encode()

	t.Run("Required_Success", func(t *testing.T) {
		t.Parallel()

		actual, err := FromQuery[string](req, Option{Key: "stringValue", RegexPattern: pattern})
		require.NoError(t, err)
		require.Equal(t, stringValue, actual)
	})

	t.Run("Required_Fail", func(t *testing.T) {
		t.Parallel()

		_, err := FromQuery[string](req, Option{Key: "stringNumberValue", RegexPattern: pattern})
		require.Error(t, err)
		require.True(t, errors.As(err, &MalformedParameterError{}))
	})
}
