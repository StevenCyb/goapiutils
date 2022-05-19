package subset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func performTest(t *testing.T, query string, testObject, expected interface{}) {
	parser := NewParser()
	actual, err := parser.Parse(query, testObject)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestSimpleMap(t *testing.T) {
	var testObject interface{} = map[string]int{"a": 1, "b": 2, "c": 3}
	performTest(t,
		"a",
		testObject,
		map[string]interface{}{})
}
