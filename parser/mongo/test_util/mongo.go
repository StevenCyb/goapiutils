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
	"go.mongodb.org/mongo-driver/bson"
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
		LogLevel:         strikememongolog.LogLevelSilent,
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
	FirstName string `bson:"first_name"`
	LatsName  string `bson:"last_name"`
	Gender    string `bson:"gender"`
	Age       int    `bson:"age"`
}

func Populate(t *testing.T, collection *mongo.Collection, items []DummyDoc) {
	itemsInterface := []interface{}{}
	for _, item := range items {
		itemsInterface = append(itemsInterface, item)
	}

	_, err := collection.InsertMany(
		context.Background(),
		itemsInterface,
	)
	require.NoError(t, err)
}

func FindCompare(t *testing.T, collection *mongo.Collection, filter interface{}, sort interface{}, items ...DummyDoc) {
	opts := &options.FindOptions{}
	if sort != nil {
		opts.Sort = sort
	}
	if filter == nil {
		filter = bson.D{}
	}

	ctx := context.TODO()
	cur, err := collection.Find(ctx, filter, opts)
	require.NoError(t, err)
	defer cur.Close(ctx)

	dbItems := []DummyDoc{}
	err = cur.All(ctx, &dbItems)
	require.NoError(t, err)

	require.Equal(t, items, dbItems)
}
