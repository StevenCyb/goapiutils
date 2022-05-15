package testutil

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/strikesecurity/strikememongo"
	"github.com/strikesecurity/strikememongo/strikememongolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoVersion = "4.2.0"

// NewStrikemongoServer creates a new strikemongo
// instance. Connection string can be obtained by
// `strikememongo.RandomDatabase()`. Keep in mind
// to stop the server after testing
// `defer mongoServer.Stop()`.
func NewStrikemongoServer(t *testing.T) *strikememongo.Server {
	options := &strikememongo.Options{
		ShouldUseReplica: false,
		LogLevel:         strikememongolog.LogLevelDebug,
		StartupTimeout:   30 * time.Second,
	}

	if mongoBin := os.Getenv("MONGO_BIN"); mongoBin != "" {
		options.MongodBin = mongoBin
	} else {
		options.MongoVersion = mongoVersion
	}

	mongoServer, err := strikememongo.StartWithOptions(options)
	require.NoError(t, err)

	return mongoServer
}

func randomID() string {
	id := uuid.New()
	return id.String()
}

// NewClientWithCollection create a new mongo client for given strikemongo
func NewClientWithCollection(t *testing.T, mongoDB *strikememongo.Server) (*mongo.Client, *mongo.Collection) {
	client, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(mongoDB.URIWithRandomDB()))
	require.NoError(t, err)

	collection := client.Database(randomID()).Collection(randomID())

	return client, collection
}

// DummyDoc is a simple dummy doc for
// mongo tests
type DummyDoc struct {
	FirstName string
	LatsName  string
	Gender    string
	Age       int
}

func Populate(t *testing.T, collection *mongo.Collection, items ...interface{}) {
	_, err := collection.InsertMany(
		context.Background(),
		items,
	)
	require.NoError(t, err)
}
