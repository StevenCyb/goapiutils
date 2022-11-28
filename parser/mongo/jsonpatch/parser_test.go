package jsonpatch

import (
	"context"
	"testing"

	"github.com/StevenCyb/goapiutils/parser/errs"
	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
	testutil "github.com/StevenCyb/goapiutils/parser/mongo/test_util"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ExecuteSuccessTest(t *testing.T, parser Parser, expect bson.A, operationSpecs ...operation.Spec) {
	t.Helper()

	actual, err := parser.Parse(operationSpecs...)
	require.NoError(t, err)
	require.Equal(t, expect, actual)
}

func ExecuteFailedTest(t *testing.T, parser Parser, expectedError error, operationSpecs ...operation.Spec) {
	t.Helper()

	_, err := parser.Parse(operationSpecs...)
	require.Equal(t, expectedError, err)
}

// DummyDoc is a simple dummy doc for mongo tests.
type DummyDoc struct {
	ID string  `bson:"_id"` //nolint:tagliatelle
	A  string  `bson:"a"`
	B  string  `bson:"b"`
	D  []int   `bson:"d"`
	C  float32 `bson:"c"`
}

func TestSingleRemoveOperation(t *testing.T) {
	t.Parallel()

	pathGroup := operation.Path("user.group")
	pathGroups := operation.Path("user.groups.3")

	ExecuteSuccessTest(t, Parser{},
		bson.A{bson.M{"$unset": "user.group"}},
		operation.Spec{
			Operation: operation.RemoveOperation,
			Path:      pathGroup,
		},
	)

	ExecuteSuccessTest(t, Parser{},
		bson.A{bson.M{"$set": bson.M{
			"user.groups": bson.M{"$concatArrays": bson.A{
				bson.M{"$slice": bson.A{"$user.groups", int64(3)}},
				bson.M{"$slice": bson.A{
					"$user.groups",
					bson.M{"$add": bson.A{1, int64(3)}},
					bson.M{"$size": "$user.groups"},
				}},
			}},
		}}},
		operation.Spec{
			Operation: operation.RemoveOperation,
			Path:      pathGroups,
		},
	)
}

func TestSingleAddOperation(t *testing.T) {
	t.Parallel()

	path := operation.Path("user.group")
	value := 1.2

	ExecuteSuccessTest(t, Parser{},
		bson.A{bson.M{"$set": bson.M{"user.group": bson.M{"$concatArrays": bson.A{"$user.group", []interface{}{1.2}}}}}},
		operation.Spec{
			Operation: operation.AddOperation,
			Path:      path,
			Value:     value,
		},
	)
}

func TestSingleReplaceOperation(t *testing.T) {
	t.Parallel()

	path := operation.Path("user.group")
	value := 1.2

	ExecuteSuccessTest(t, Parser{},
		bson.A{bson.M{"$set": bson.M{"user.group": 1.2}}},
		operation.Spec{
			Operation: operation.ReplaceOperation,
			Path:      path,
			Value:     value,
		},
	)
}

func TestSingleMoveOperation(t *testing.T) {
	t.Parallel()

	path := operation.Path("user.a")
	from := operation.Path("user.a_tmp")
	value := 1.2

	ExecuteSuccessTest(t, Parser{},
		bson.A{
			bson.M{"$set": bson.M{"user.a": "$user.a_tmp"}},
			bson.M{"$unset": "user.a_tmp"},
		},
		operation.Spec{
			Operation: operation.MoveOperation,
			Path:      path,
			Value:     value,
			From:      from,
		},
	)
}

func TestSingleCopyOperation(t *testing.T) {
	t.Parallel()

	path := operation.Path("user.a")
	from := operation.Path("user.a_tmp")
	value := 1.2

	ExecuteSuccessTest(t, Parser{},
		bson.A{bson.M{"$set": bson.M{"user.a": "$user.a_tmp"}}},
		operation.Spec{
			Operation: operation.CopyOperation,
			Path:      path,
			Value:     value,
			From:      from,
		},
	)
}

func TestInvalidOperation(t *testing.T) {
	t.Parallel()

	name := "MustFailPolicy"
	path := operation.Path("user.a")

	ExecuteFailedTest(t,
		Parser{
			[]Policy{DisallowPathPolicy{Details: name, Path: path}},
		},
		errs.NewErrUnexpectedInput(operation.Spec{}),
		operation.Spec{},
	)
}

func TestPolicyViolation(t *testing.T) {
	t.Parallel()

	name := "MustFailPolicy"
	path := operation.Path("user.a")

	ExecuteFailedTest(t,
		Parser{
			[]Policy{DisallowPathPolicy{Details: name, Path: path}},
		},
		errs.NewErrPolicyViolation(name),
		operation.Spec{Operation: operation.RemoveOperation, Path: path},
	)
}

func TestInterpretation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	server := testutil.NewStrikemongoServer(t)
	mongoClient, collection, database := testutil.NewClientWithCollection(t, server)

	//nolint:errcheck
	t.Cleanup(func() {
		server.Stop()
		mongoClient.Disconnect(ctx)
		database.Drop(ctx)
	})

	items := []DummyDoc{
		{ID: "1", A: "a1", B: "b1", C: float32(11), D: []int{1}},
		{ID: "2", A: "a2", B: "b2", C: float32(12), D: []int{2}},
		{ID: "3", A: "a3", B: "b3", C: float32(13), D: []int{3}},
		{ID: "4", A: "a4", B: "b4", C: float32(14), D: []int{4}},
		{ID: "5", A: "a5", B: "b5", C: float32(15), D: []int{5}},
	}
	itemsInterface := []interface{}{}

	for _, item := range items {
		itemsInterface = append(itemsInterface, item)
	}

	testutil.Populate(t, collection, itemsInterface)

	testRemoveOperation(t, collection, items[0])
	testAddOperation(t, collection, items[1])
	testReplaceOperation(t, collection, items[2])
	testMoveOperation(t, collection, items[3])
	testCopyOperation(t, collection, items[4])
}

func testRemoveOperation(t *testing.T, collection *mongo.Collection, item DummyDoc) {
	t.Helper()

	ctx := context.Background()
	parser := Parser{}
	after := options.After
	updateOptions := &options.FindOneAndUpdateOptions{ReturnDocument: &after}
	query, err := parser.Parse(operation.Spec{Operation: operation.RemoveOperation, Path: "a"})
	require.NoError(t, err)

	item.A = ""

	result := collection.FindOneAndUpdate(ctx, bson.D{{Key: "_id", Value: item.ID}}, query, updateOptions)
	require.NoError(t, result.Err())

	var resultingDocument DummyDoc
	err = result.Decode(&resultingDocument)
	require.NoError(t, err)
	require.Equal(t, item, resultingDocument)
}

func testAddOperation(t *testing.T, collection *mongo.Collection, item DummyDoc) {
	t.Helper()

	ctx := context.Background()
	parser := Parser{}
	after := options.After
	updateOptions := &options.FindOneAndUpdateOptions{ReturnDocument: &after}
	value := 3
	query, err := parser.Parse(operation.Spec{Operation: operation.AddOperation, Path: "d", Value: value})
	require.NoError(t, err)

	item.D = append(item.D, value)

	result := collection.FindOneAndUpdate(ctx, bson.D{{Key: "_id", Value: item.ID}}, query, updateOptions)
	require.NoError(t, result.Err())

	var resultingDocument DummyDoc
	err = result.Decode(&resultingDocument)
	require.NoError(t, err)
	require.Equal(t, item, resultingDocument)
}

func testReplaceOperation(t *testing.T, collection *mongo.Collection, item DummyDoc) {
	t.Helper()

	ctx := context.Background()
	parser := Parser{}
	after := options.After
	updateOptions := &options.FindOneAndUpdateOptions{ReturnDocument: &after}
	value := float32(99.99)
	query, err := parser.Parse(operation.Spec{Operation: operation.ReplaceOperation, Path: "c", Value: value})
	require.NoError(t, err)

	item.C = 99.99

	result := collection.FindOneAndUpdate(ctx, bson.D{{Key: "_id", Value: item.ID}}, query, updateOptions)
	require.NoError(t, result.Err())

	var resultingDocument DummyDoc
	err = result.Decode(&resultingDocument)
	require.NoError(t, err)
	require.Equal(t, item, resultingDocument)
}

func testMoveOperation(t *testing.T, collection *mongo.Collection, item DummyDoc) {
	t.Helper()

	ctx := context.Background()
	parser := Parser{}
	after := options.After
	updateOptions := &options.FindOneAndUpdateOptions{ReturnDocument: &after}
	query, err := parser.Parse(operation.Spec{Operation: operation.MoveOperation, Path: "a", From: "b"})
	require.NoError(t, err)

	item.A = item.B
	item.B = ""

	result := collection.FindOneAndUpdate(ctx, bson.D{{Key: "_id", Value: item.ID}}, query, updateOptions)
	require.NoError(t, result.Err())

	var resultingDocument DummyDoc
	err = result.Decode(&resultingDocument)
	require.NoError(t, err)
	require.Equal(t, item, resultingDocument)
}

func testCopyOperation(t *testing.T, collection *mongo.Collection, item DummyDoc) {
	t.Helper()

	ctx := context.Background()
	parser := Parser{}
	after := options.After
	updateOptions := &options.FindOneAndUpdateOptions{ReturnDocument: &after}
	query, err := parser.Parse(operation.Spec{Operation: operation.CopyOperation, Path: "a", From: "b"})
	require.NoError(t, err)

	item.A = item.B

	result := collection.FindOneAndUpdate(ctx, bson.D{{Key: "_id", Value: item.ID}}, query, updateOptions)
	require.NoError(t, result.Err())

	var resultingDocument DummyDoc
	err = result.Decode(&resultingDocument)
	require.NoError(t, err)
	require.Equal(t, item, resultingDocument)
}
