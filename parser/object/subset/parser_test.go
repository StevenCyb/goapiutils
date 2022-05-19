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

	t.Run("ExistingField_Success", func(t *testing.T) {
		performTest(t,
			"a=extracted_a",
			testObject,
			map[string]interface{}{"extracted_a": 1})
	})

	t.Run("NonExistingField_Success", func(t *testing.T) {
		performTest(t,
			"x=x",
			testObject,
			map[string]interface{}{})
	})
}

func TestMapWithArray(t *testing.T) {
	var testObject interface{} = map[string]interface{}{
		"nested": []map[string]string{
			{"a": "A", "b": "B", "c": "C"},
			{"d": "D", "e": "E", "f": "F"},
		},
	}

	t.Run("ExistingField_Success", func(t *testing.T) {
		performTest(t,
			"nested=arr",
			testObject,
			map[string]interface{}{"arr": []map[string]string{
				{"a": "A", "b": "B", "c": "C"},
				{"d": "D", "e": "E", "f": "F"},
			}})
	})

	t.Run("NonExistingField_Success", func(t *testing.T) {
		performTest(t,
			"nested.a=no",
			testObject,
			map[string]interface{}{})
	})
}

func TestNestedMaps(t *testing.T) {
	var testObject interface{} = map[string]interface{}{
		"nested": map[string]map[string]string{
			"german": {
				"greet": "hallo",
			},
			"english": {
				"greet": "hello",
			},
		},
	}

	t.Run("ExistingField_Success", func(t *testing.T) {
		performTest(t,
			"nested.english.greet=greet",
			testObject,
			map[string]interface{}{"greet": "hello"})
	})

	t.Run("NonExistingField_Success", func(t *testing.T) {
		performTest(t,
			"nested.english.3=no",
			testObject,
			map[string]interface{}{})
	})
}

func TestCombinedSubsetsMap(t *testing.T) {
	var testObject interface{} = map[string]int{"a": 1, "b": 2, "c": 3}

	t.Run("ExistingField_Success", func(t *testing.T) {
		performTest(t,
			"a=extracted_a,b=extracted_b",
			testObject,
			map[string]interface{}{"extracted_a": 1, "extracted_b": 2})
	})

	t.Run("WithOneNonExistingField_Success", func(t *testing.T) {
		performTest(t,
			"a=extracted_a,x=x",
			testObject,
			map[string]interface{}{"extracted_a": 1})
	})
}
