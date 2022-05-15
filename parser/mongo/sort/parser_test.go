package sort

import (
	"testing"

	"github.com/StevenCyb/goquery/errs"
	testutil "github.com/StevenCyb/goquery/parser/mongo/test_util"
	"github.com/StevenCyb/goquery/tokenizer"

	"go.mongodb.org/mongo-driver/bson"
)

func TestParsing(t *testing.T) {
	t.Run("Query", func(t *testing.T) {
		t.Run("WithEmptyQuery_Success", func(t *testing.T) {
			testutil.ExecuteSuccessTest(t,
				NewParser(nil),
				"",
				bson.D{},
			)
		})

		t.Run("WithSingleAscCriteria_Success", func(t *testing.T) {
			testutil.ExecuteSuccessTest(t,
				NewParser(nil),
				"name=asc",
				bson.D{bson.E{Key: "name", Value: 1}},
			)
			testutil.ExecuteSuccessTest(t,
				NewParser(nil),
				"name=1",
				bson.D{bson.E{Key: "name", Value: 1}},
			)
		})

		t.Run("WithSingleDescCriteria_Success", func(t *testing.T) {
			testutil.ExecuteSuccessTest(t,
				NewParser(nil),
				"name=desc",
				bson.D{bson.E{Key: "name", Value: -1}},
			)
			testutil.ExecuteSuccessTest(t,
				NewParser(nil),
				"name=-1",
				bson.D{bson.E{Key: "name", Value: -1}},
			)
		})

		t.Run("WithMultipleSortingCriteria_Success", func(t *testing.T) {
			testutil.ExecuteSuccessTest(t,
				NewParser(nil),
				"firstName=asc,age=desc, something = asc",
				bson.D{
					bson.E{Key: "firstName", Value: 1},
					bson.E{Key: "age", Value: -1},
					bson.E{Key: "something ", Value: 1},
				},
			)
		})

		t.Run("WithUnknownLiteral_Fail", func(t *testing.T) {
			testutil.ExecuteFailedTest(t,
				NewParser(nil),
				"firstName=?",
				errs.NewErrUnexpectedTokenType(11, "FIELD_NAME", "SORT_CRITERIA"),
			)
		})
		t.Run("WithUnexpectedSeparator_Fail", func(t *testing.T) {
			testutil.ExecuteFailedTest(t,
				NewParser(nil),
				"firstName=asc+lastName=asc",
				errs.NewErrUnexpectedTokenType(22, "FIELD_NAME", ","),
			)
		})
	})

	t.Run("WithPolicy", func(t *testing.T) {
		t.Run("WithAllowedFieldnames_Success", func(t *testing.T) {
			testutil.ExecuteSuccessTest(t,
				NewParser(
					tokenizer.NewPolicy(tokenizer.WHITELIST_POLICY, "a", "b")),
				"a=asc,b=asc",
				bson.D{bson.E{Key: "a", Value: 1}, bson.E{Key: "b", Value: 1}},
			)
		})

		t.Run("WithDisallowedFieldnames_Fail", func(t *testing.T) {
			testutil.ExecuteFailedTest(t,
				NewParser(
					tokenizer.NewPolicy(tokenizer.WHITELIST_POLICY, "a", "b")),
				"a=asc,b=desc,c=desc",
				errs.NewErrPolicyViolation("c"),
			)
		})
	})
}

func TestInterpretation(t *testing.T) {
	server := testutil.NewStrikemongoServer(t)
	defer server.Stop()
	// TODO user strikemongo for testing
}
