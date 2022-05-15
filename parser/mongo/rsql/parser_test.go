package rsql

import (
	"context"
	"regexp"
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

		t.Run("WithDifferentLiterals", func(t *testing.T) {
			t.Run("==STRING_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`firstName=="steven"`,
					bson.D{bson.E{Key: "firstName", Value: "steven"}},
				)
			})

			t.Run("==INT_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`year==2022`,
					bson.D{bson.E{Key: "year", Value: int64(2022)}},
				)
			})

			t.Run("==FLOAT_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`pi==3.14159265`,
					bson.D{bson.E{Key: "pi", Value: float64(3.14159265)}},
				)
			})

			t.Run("==BOOL_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`is==TRUE`,
					bson.D{bson.E{Key: "is", Value: true}},
				)
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`is==false`,
					bson.D{bson.E{Key: "is", Value: false}},
				)
			})

			t.Run("=in=SLICE_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`coll=in=(1, "a", "b",2)`,
					bson.D{bson.E{Key: "coll", Value: bson.E{
						Key:   "$in",
						Value: bson.A{int64(1), "a", "b", int64(2)}}}},
				)
			})
		})

		t.Run("WithSingleComparisonOperation", func(t *testing.T) {
			t.Run("==_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`firstName=="steven"`,
					bson.D{bson.E{Key: "firstName", Value: "steven"}},
				)
			})

			t.Run("!=_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`x!=10`,
					bson.D{bson.E{Key: "x",
						Value: bson.E{Key: "$ne", Value: int64(10)}}},
				)
			})

			t.Run("=sw=_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`msg=sw="LOG_"`,
					bson.D{bson.E{Key: "msg", Value: *regexp.MustCompile("^LOG_")}},
				)
			})

			t.Run("=ew=_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`word=ew="ed"`,
					bson.D{bson.E{Key: "word", Value: *regexp.MustCompile("ed$")}},
				)
			})

			t.Run("=gt=_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`x=gt=10`,
					bson.D{bson.E{Key: "x",
						Value: bson.E{Key: "$gt", Value: int64(10)}}},
				)
			})

			t.Run("=ge=_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`x=ge=10`,
					bson.D{bson.E{Key: "x",
						Value: bson.E{Key: "$gte", Value: int64(10)}}},
				)
			})

			t.Run("=lt=_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`x=lt=10`,
					bson.D{bson.E{Key: "x",
						Value: bson.E{Key: "$lt", Value: int64(10)}}},
				)
			})

			t.Run("=le=_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`x=le=10`,
					bson.D{bson.E{Key: "x",
						Value: bson.E{Key: "$lte", Value: int64(10)}}},
				)
			})

			t.Run("=in=_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`coll=in=(1, "a", "b",2)`,
					bson.D{bson.E{Key: "coll", Value: bson.E{
						Key:   "$in",
						Value: bson.A{int64(1), "a", "b", int64(2)}}}},
				)
			})

			t.Run("=out=_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`coll=out=(1, "a", "b",2)`,
					bson.D{bson.E{Key: "coll", Value: bson.E{
						Key:   "$nin",
						Value: bson.A{int64(1), "a", "b", int64(2)}}}},
				)
			})
		})

		t.Run("WithMultipleComparisonOperation", func(t *testing.T) {
			t.Run("WithSingleAnd_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`firstName=="steven";age=ge=18`,
					bson.D{
						bson.E{
							Key: "$and", Value: bson.A{
								bson.E{Key: "firstName", Value: "steven"},
								bson.E{Key: "age",
									Value: bson.E{Key: "$gte", Value: int64(18)}},
							}}},
				)
			})

			t.Run("WithSingleOr_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`level=="error",level=="warning"`,
					bson.D{
						bson.E{
							Key: "$or", Value: bson.A{
								bson.E{Key: "level", Value: "error"},
								bson.E{Key: "level", Value: "warning"},
							}}},
				)
			})

			t.Run("WithAndChain_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`firstName=="steven";age=ge=18;gender=="male"`,
					bson.D{
						bson.E{
							Key: "$and", Value: bson.A{
								bson.E{Key: "firstName", Value: "steven"},
								bson.E{Key: "age",
									Value: bson.E{Key: "$gte", Value: int64(18)}},
								bson.E{Key: "gender", Value: "male"},
							}}},
				)
			})

			t.Run("WithOrChain_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`level=="panic",level=="error",level=="warning"`,
					bson.D{
						bson.E{
							Key: "$or", Value: bson.A{
								bson.E{Key: "level", Value: "panic"},
								bson.E{Key: "level", Value: "error"},
								bson.E{Key: "level", Value: "warning"},
							}}},
				)
			})

			t.Run("WithMixed_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`a==1,a==2,a==3,b==1;c==1`,
					bson.D{
						bson.E{Key: "$or", Value: bson.A{
							bson.E{Key: "a", Value: int64(1)},
							bson.E{Key: "a", Value: int64(2)},
							bson.E{Key: "a", Value: int64(3)},
							bson.E{Key: "$and", Value: bson.A{
								bson.E{Key: "b", Value: int64(1)},
								bson.E{Key: "c", Value: int64(1)}}}}}},
				)

				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`a==1;b==1,a==2;b==2`,
					bson.D{
						bson.E{Key: "$and", Value: bson.A{
							bson.E{Key: "a", Value: int64(1)},
							bson.E{Key: "$or", Value: bson.A{
								bson.E{Key: "b", Value: int64(1)},
								bson.E{Key: "$and", Value: bson.A{
									bson.E{Key: "a", Value: int64(2)},
									bson.E{Key: "b", Value: int64(2)}}}}}}}},
				)
			})

			t.Run("WithContexted_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`(a==1;b==1),(a==2;b==2),(a==3;b==3)`,
					bson.D{
						bson.E{Key: "$or", Value: bson.A{
							bson.E{Key: "$and", Value: bson.A{
								bson.E{Key: "a", Value: int64(1)},
								bson.E{Key: "b", Value: int64(1)},
							}},
							bson.E{Key: "$and", Value: bson.A{
								bson.E{Key: "a", Value: int64(2)},
								bson.E{Key: "b", Value: int64(2)},
							}},
							bson.E{Key: "$and", Value: bson.A{
								bson.E{Key: "a", Value: int64(3)},
								bson.E{Key: "b", Value: int64(3)},
							}},
						}},
					},
				)
			})

			t.Run("WithNestedContext_Success", func(t *testing.T) {
				testutil.ExecuteSuccessTest(t,
					NewParser(nil),
					`(a==1;b==1),((a==2,b==2);(a==3,b==3))`,
					bson.D{
						bson.E{Key: "$or", Value: bson.A{
							bson.E{Key: "$and", Value: bson.A{
								bson.E{Key: "a", Value: int64(1)},
								bson.E{Key: "b", Value: int64(1)},
							}},
							bson.E{Key: "$and", Value: bson.A{
								bson.E{Key: "$or", Value: bson.A{
									bson.E{Key: "a", Value: int64(2)},
									bson.E{Key: "b", Value: int64(2)},
								}},
								bson.E{Key: "$or", Value: bson.A{
									bson.E{Key: "a", Value: int64(3)},
									bson.E{Key: "b", Value: int64(3)},
								}},
							}},
						}},
					},
				)
			})
		})

		t.Run("WithUnknownCompareOperation_Fail", func(t *testing.T) {
			testutil.ExecuteFailedTest(t,
				NewParser(nil),
				`x=7`,
				errs.NewErrUnexpectedToken(1, "="),
			)
		})

		t.Run("WithIncompleteExpression_Fail", func(t *testing.T) {
			testutil.ExecuteFailedTest(t,
				NewParser(nil),
				`x==7;`,
				errs.NewErrUnexpectedInputEnd("FIELD_NAME"),
			)
		})

		t.Run("WithWrongCompareOperation_Fail", func(t *testing.T) {
			testutil.ExecuteFailedTest(t,
				NewParser(nil),
				`x==(1,2,3)`,
				errs.NewErrUnexpectedTokenType(3, "(", "LITERAL"),
			)
			testutil.ExecuteFailedTest(t,
				NewParser(nil),
				`x=in=3`,
				errs.NewErrUnexpectedTokenType(6, "NUMERIC_LITERAL", "("),
			)
		})

		t.Run("WithNotClosedContext_Fail", func(t *testing.T) {
			testutil.ExecuteFailedTest(t,
				NewParser(nil),
				`(x==7`,
				errs.NewErrUnexpectedInputEnd(")"),
			)
		})
	})

	t.Run("WithPolicy", func(t *testing.T) {
		t.Run("WithAllowedFieldnames_Success", func(t *testing.T) {
			testutil.ExecuteSuccessTest(t,
				NewParser(
					tokenizer.NewPolicy(tokenizer.WHITELIST_POLICY, "name", "age")),
				`name=="steven",age=ge=18`,
				bson.D{
					bson.E{Key: "$or", Value: bson.A{
						bson.E{Key: "name", Value: "steven"},
						bson.E{Key: "age", Value: bson.E{Key: "$gte", Value: int64(18)}}}}},
			)
		})

		t.Run("WithDisallowedFieldnames_Fail", func(t *testing.T) {
			testutil.ExecuteFailedTest(t,
				NewParser(
					tokenizer.NewPolicy(tokenizer.WHITELIST_POLICY, "name", "age")),
				`name=="steven",age=ge=18,gender="male"`,
				errs.NewErrPolicyViolation("gender"),
			)
		})
	})
}

func TestInterpretation(t *testing.T) {
	server := testutil.NewStrikemongoServer(t)
	defer server.Stop()
	mongoClient, collection := testutil.NewClientWithCollection(t, server)
	defer mongoClient.Disconnect(context.Background())
	testutil.Populate(t, collection,
		testutil.DummyDoc{FirstName: "Max", LatsName: "Muster", Gender: "male", Age: 52},
		testutil.DummyDoc{FirstName: "Alexa", LatsName: "Amaizon", Gender: "female", Age: 22},
		testutil.DummyDoc{FirstName: "Tina", LatsName: "Someone", Gender: "female", Age: 33},
		testutil.DummyDoc{FirstName: "Samal", LatsName: "Someone", Gender: "male", Age: 26},
	)
	// TODO user strikemongo for testing
}
