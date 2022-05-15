package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/strikesecurity/strikememongo"
)

const strikememongoVersion = "4.2.0"

// NewStrikemongoServer creates a new strikemongo
// instance. Connection string can be obtained by
// `strikememongo.RandomDatabase()`. Keep in mind
// to stop the server after testing
// `defer mongoServer.Stop()`.
func NewStrikemongoServer(t *testing.T) *strikememongo.Server {
	mongoServer, err := strikememongo.Start(strikememongoVersion)
	require.NoError(t, err)

	return mongoServer
}
