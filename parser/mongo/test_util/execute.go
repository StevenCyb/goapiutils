package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

type Parser interface {
	Parse(query string) (bson.D, error)
}

func ExecuteSuccessTest(t *testing.T, parser Parser, query string, expect bson.D) {
	actual, err := parser.Parse(query)
	require.NoError(t, err)
	require.Equal(t, expect, actual)
}

func ExecuteFailedTest(t *testing.T, parser Parser, query string, expectedError error) {
	_, err := parser.Parse(query)
	require.Equal(t, expectedError, err)
}
