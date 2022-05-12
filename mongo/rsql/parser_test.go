package rsql

import (
	"regexp"
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

func TestEmpty(t *testing.T) {
	executeSuccessTest(t,
		"",
		bson.D{},
	)
}

func TestLiterals(t *testing.T) {
	t.Run("==STRING", func(t *testing.T) {
		executeSuccessTest(t,
			`firstName=="steven"`,
			bson.D{bson.E{Key: "firstName", Value: "steven"}},
		)
	})

	t.Run("==INT", func(t *testing.T) {
		executeSuccessTest(t,
			`year==2022`,
			bson.D{bson.E{Key: "year", Value: int64(2022)}},
		)
	})

	t.Run("==FLOAT", func(t *testing.T) {
		executeSuccessTest(t,
			`pi==3.14159265`,
			bson.D{bson.E{Key: "pi", Value: float64(3.14159265)}},
		)
	})

	t.Run("==BOOL", func(t *testing.T) {
		executeSuccessTest(t,
			`is==TRUE`,
			bson.D{bson.E{Key: "is", Value: true}},
		)
		executeSuccessTest(t,
			`is==false`,
			bson.D{bson.E{Key: "is", Value: false}},
		)
	})

	t.Run("=in=SLICE", func(t *testing.T) {
		executeSuccessTest(t,
			`coll=in=(1, "a", "b",2)`,
			bson.D{bson.E{Key: "coll", Value: bson.E{
				Key:   "$in",
				Value: bson.A{int64(1), "a", "b", int64(2)}}}},
		)
	})
}

func TestSingleExpression(t *testing.T) {
	t.Run("==", func(t *testing.T) {
		executeSuccessTest(t,
			`firstName=="steven"`,
			bson.D{bson.E{Key: "firstName", Value: "steven"}},
		)
	})

	t.Run("!=", func(t *testing.T) {
		executeSuccessTest(t,
			`x!=10`,
			bson.D{bson.E{Key: "x",
				Value: bson.E{Key: "$ne", Value: int64(10)}}},
		)
	})

	t.Run("=sw=", func(t *testing.T) {
		executeSuccessTest(t,
			`msg=sw="LOG_"`,
			bson.D{bson.E{Key: "msg", Value: *regexp.MustCompile("^LOG_")}},
		)
	})

	t.Run("=ew=", func(t *testing.T) {
		executeSuccessTest(t,
			`word=ew="ed"`,
			bson.D{bson.E{Key: "word", Value: *regexp.MustCompile("ed$")}},
		)
	})

	t.Run("=gt=", func(t *testing.T) {
		executeSuccessTest(t,
			`x=gt=10`,
			bson.D{bson.E{Key: "x",
				Value: bson.E{Key: "$gt", Value: int64(10)}}},
		)
	})

	t.Run("=ge=", func(t *testing.T) {
		executeSuccessTest(t,
			`x=ge=10`,
			bson.D{bson.E{Key: "x",
				Value: bson.E{Key: "$gte", Value: int64(10)}}},
		)
	})

	t.Run("=lt=", func(t *testing.T) {
		executeSuccessTest(t,
			`x=lt=10`,
			bson.D{bson.E{Key: "x",
				Value: bson.E{Key: "$lt", Value: int64(10)}}},
		)
	})

	t.Run("=le=", func(t *testing.T) {
		executeSuccessTest(t,
			`x=le=10`,
			bson.D{bson.E{Key: "x",
				Value: bson.E{Key: "$lte", Value: int64(10)}}},
		)
	})

	t.Run("=in=", func(t *testing.T) {
		executeSuccessTest(t,
			`coll=in=(1, "a", "b",2)`,
			bson.D{bson.E{Key: "coll", Value: bson.E{
				Key:   "$in",
				Value: bson.A{int64(1), "a", "b", int64(2)}}}},
		)
	})

	t.Run("=out=", func(t *testing.T) {
		executeSuccessTest(t,
			`coll=out=(1, "a", "b",2)`,
			bson.D{bson.E{Key: "coll", Value: bson.E{
				Key:   "$nin",
				Value: bson.A{int64(1), "a", "b", int64(2)}}}},
		)
	})
}

func TestMultipleCompositedExpressions(t *testing.T) {
	t.Run("SingleAnd", func(t *testing.T) {
		executeSuccessTest(t,
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

	t.Run("SingleOr", func(t *testing.T) {
		executeSuccessTest(t,
			`level=="error",level=="warning"`,
			bson.D{
				bson.E{
					Key: "$or", Value: bson.A{
						bson.E{Key: "level", Value: "error"},
						bson.E{Key: "level", Value: "warning"},
					}}},
		)
	})

	t.Run("AndChain", func(t *testing.T) {
		executeSuccessTest(t,
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

	t.Run("OrChain", func(t *testing.T) {
		executeSuccessTest(t,
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

	t.Run("Mixed", func(t *testing.T) {
		executeSuccessTest(t,
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

		executeSuccessTest(t,
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

	t.Run("Contexted", func(t *testing.T) {
		executeSuccessTest(t,
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

	t.Run("NestedContext", func(t *testing.T) {
		executeSuccessTest(t,
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
}

func TestMalformedRsql(t *testing.T) {
	t.Run("UnknownCompareOperation", func(t *testing.T) {
		executeFailedTest(t,
			`x=7`,
			errs.NewErrUnexpectedToken(1, "="),
		)
	})

	t.Run("IncompleteExpression", func(t *testing.T) {
		executeFailedTest(t,
			`x==7;`,
			errs.NewErrUnexpectedInputEnd("FIELD_NAME"),
		)
	})

	t.Run("WrongCompareOperation", func(t *testing.T) {
		executeFailedTest(t,
			`x==(1,2,3)`,
			errs.NewErrUnexpectedTokenType(3, "(", "LITERAL"),
		)
		executeFailedTest(t,
			`x=in=3`,
			errs.NewErrUnexpectedTokenType(6, "NUMERIC_LITERAL", "("),
		)
	})

	t.Run("NotClosedContext", func(t *testing.T) {
		executeFailedTest(t,
			`(x==7`,
			errs.NewErrUnexpectedInputEnd(")"),
		)
	})
}

func TestPolicy(t *testing.T) {
	t.Run("AllowByPolicy", func(t *testing.T) {
		parser := NewParser(
			tokenizer.NewPolicy(tokenizer.WHITELIST_POLICY, "name", "age"))
		_, err := parser.Parse(`name=="steven",age=ge=18`)
		require.NoError(t, err)
	})

	t.Run("DisallowedByPolicy", func(t *testing.T) {
		parser := NewParser(
			tokenizer.NewPolicy(tokenizer.WHITELIST_POLICY, "name", "age"))
		_, err := parser.Parse(`name=="steven",age=ge=18,gender="male"`)
		require.Equal(t, errs.NewErrPolicyViolation("gender"), err)
	})
}
