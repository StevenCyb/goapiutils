package sort

import (
	"testing"

	"github.com/StevenCyb/goquery/errs"
	"github.com/StevenCyb/goquery/tokenizer"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func executeSuccessTest(t *testing.T, query string, expect bson.D) {
	parser := NewParser(nil)
	actual, err := parser.Parse(query)
	require.NoError(t, err)
	require.Equal(t, expect, actual)
}

func executeFailedTest(t *testing.T, query string, expectedError error) {
	parser := NewParser(nil)
	_, err := parser.Parse(query)
	require.Equal(t, expectedError, err)
}

func TestValidExamples(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		executeSuccessTest(t,
			"",
			bson.D{},
		)
	})

	t.Run("SingleAsc", func(t *testing.T) {
		executeSuccessTest(t,
			"name=asc",
			bson.D{bson.E{Key: "name", Value: 1}},
		)
		executeSuccessTest(t,
			"name=1",
			bson.D{bson.E{Key: "name", Value: 1}},
		)
	})

	t.Run("SingleDesc", func(t *testing.T) {
		executeSuccessTest(t,
			"name=desc",
			bson.D{bson.E{Key: "name", Value: -1}},
		)
		executeSuccessTest(t,
			"name=-1",
			bson.D{bson.E{Key: "name", Value: -1}},
		)
	})

	t.Run("Multiple", func(t *testing.T) {
		executeSuccessTest(t,
			"firstName=asc,age=desc, something = asc",
			bson.D{
				bson.E{Key: "firstName", Value: 1},
				bson.E{Key: "age", Value: -1},
				bson.E{Key: "something ", Value: 1},
			},
		)
	})
}

func TestMalformedExpression(t *testing.T) {
	t.Run("UnknownLiteral", func(t *testing.T) {
		executeFailedTest(t,
			"firstName=?",
			errs.NewErrUnexpectedTokenType(11, "FIELD_NAME", "SORT_CRITERIA"),
		)
	})
	t.Run("UnexpectedSeparator", func(t *testing.T) {
		executeFailedTest(t,
			"firstName=asc+lastName=asc",
			errs.NewErrUnexpectedTokenType(22, "FIELD_NAME", ","),
		)
	})
}

func TestPolicy(t *testing.T) {
	t.Run("AllowByPolicy", func(t *testing.T) {
		parser := NewParser(
			tokenizer.NewPolicy(tokenizer.WHITELIST_POLICY, "a", "b"))
		_, err := parser.Parse("a=asc,b=asc")
		require.NoError(t, err)
	})

	t.Run("DisallowedByPolicy", func(t *testing.T) {
		parser := NewParser(
			tokenizer.NewPolicy(tokenizer.WHITELIST_POLICY, "a", "b"))
		_, err := parser.Parse("a=asc,b=desc,c=desc")
		require.Equal(t, errs.NewErrPolicyViolation("c"), err)
	})
}
