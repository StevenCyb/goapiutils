//nolint:dupl,ireturn
package validator

import (
	"errors"
	"reflect"
	"regexp"
	"sort"
	"testing"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/rule"
	"github.com/stretchr/testify/require"
)

type demoRule struct{}

func (d demoRule) NewInstance(patch string, kind reflect.Kind, instance interface{}, value string) (rule.Rule, error) { //nolint:ireturn,lll
	return nil, nil
}

func (d demoRule) NewInheritInstance(patch string, kind reflect.Kind, instance interface{}) (rule.Rule, error) { //nolint:ireturn,lll
	return nil, nil
}

func (d demoRule) Validate(operationSpec operation.Spec) error {
	return nil
}

func compareMapWithRandomArray[T any](t *testing.T, expected, actual map[operation.Path][]T) {
	t.Helper()

	for expectedPath, expectedUnsortedArray := range expected {
		expectedArr := expectedUnsortedArray

		sort.SliceStable(expectedArr, func(i, j int) bool {
			iKind := reflect.TypeOf(expectedArr[i]).String()
			jKind := reflect.TypeOf(expectedArr[j]).String()

			return iKind < jKind
		})

		if actualUnsortedArray, exists := actual[expectedPath]; exists {
			actualArr := actualUnsortedArray

			sort.SliceStable(actualArr, func(i, j int) bool {
				iKind := reflect.TypeOf(actualArr[i]).String()
				jKind := reflect.TypeOf(actualArr[j]).String()

				return iKind < jKind
			})

			require.Equal(t, expectedArr, actualArr)
		} else {
			t.Fatalf("missing rule for '%s'", expectedPath)
		}
	}
}

func TestInstantiation(t *testing.T) {
	t.Parallel()

	t.Run("ValidReference_Success", func(t *testing.T) {
		t.Parallel()

		validator, err := NewValidator(reflect.TypeOf(struct{}{}))
		require.NoError(t, err)
		require.NotNil(t, validator)
	})

	t.Run("InvalidReference_Fail", func(t *testing.T) {
		t.Parallel()

		validator, err := NewValidator(reflect.TypeOf(nil))
		require.Error(t, err)
		require.True(t, errors.Is(err, ErrReferenceIsNil))
		require.Nil(t, validator)
	})
}

func TestRegisterRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidRule_Success", func(t *testing.T) {
		t.Parallel()

		err = validator.RegisterRule("jp_name", demoRule{})
		require.NoError(t, err)
	})

	t.Run("DuplicateRuleName_Fail", func(t *testing.T) {
		t.Parallel()

		err = validator.RegisterRule("jp_disallow", demoRule{})
		require.Error(t, err)
		require.Equal(t, ErrDuplicateRuleTags, err)
	})

	t.Run("MissingPrefix_Fail", func(t *testing.T) {
		t.Parallel()

		err = validator.RegisterRule("name", demoRule{})
		require.Error(t, err)
		require.Equal(t, ErrMissingPrefix, err)
	})
}

func TestUseReferenceWithSimpleTypes(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	err = validator.UseReference(reflect.TypeOf(struct { //nolint:govet
		A string
		B string            `bson:"b"`
		C int               `bson:"c"`
		D []string          `bson:"d"`
		E []int             `bson:"e"`
		F map[string]string `bson:"f"`
		G map[string]int    `bson:"g"`
		H struct {
			A string
			B string `bson:"b"`
		} `bson:"h"`
	}{}))
	expectedRule := map[operation.Path][]rule.Rule{
		"b": {&rule.MatchingKindRule{Instance: ""}},
		"c": {&rule.MatchingKindRule{Instance: 0}},
		"d": {&rule.MatchingKindRule{Instance: []string{}}},
		"e": {&rule.MatchingKindRule{Instance: []int{}}},
		"f": {&rule.MatchingKindRule{Instance: map[string]string{}}},
		"g": {&rule.MatchingKindRule{Instance: map[string]int{}}},
		"h": {&rule.MatchingKindRule{Instance: struct {
			A string
			B string `bson:"b"`
		}{}}},
		"h.b": {&rule.MatchingKindRule{Instance: ""}},
	}
	expectedWildcardRules := map[operation.Path][]rule.Rule{
		"d.*": {&rule.MatchingKindRule{Instance: ""}},
		"e.*": {&rule.MatchingKindRule{Instance: 0}},
		"f.*": {&rule.MatchingKindRule{Instance: ""}},
		"g.*": {&rule.MatchingKindRule{Instance: 0}},
	}

	require.NoError(t, err)
	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceComplexStruct(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	err = validator.UseReference(reflect.TypeOf(struct {
		A struct {
			A struct {
				B string `bson:"b"`
			} `bson:"a"`
		} `bson:"a"`
	}{}))
	expectedRule := map[operation.Path][]rule.Rule{
		"a": {&rule.MatchingKindRule{Instance: struct {
			A struct {
				B string `bson:"b"`
			} `bson:"a"`
		}{}}},
		"a.a": {&rule.MatchingKindRule{Instance: struct {
			B string `bson:"b"`
		}{}}},
		"a.a.b": {&rule.MatchingKindRule{Instance: ""}},
	}
	expectedWildcardRules := map[operation.Path][]rule.Rule{}

	require.NoError(t, err)
	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceComplexArray(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	err = validator.UseReference(reflect.TypeOf(struct {
		B []struct {
			A string `bson:"a"`
		} `bson:"b"`
	}{}))
	expectedRule := map[operation.Path][]rule.Rule{
		"b": {&rule.MatchingKindRule{Instance: []struct {
			A string `bson:"a"`
		}{}}},
	}
	expectedWildcardRules := map[operation.Path][]rule.Rule{
		"b.*": {&rule.MatchingKindRule{Instance: struct {
			A string `bson:"a"`
		}{}}},
		"b.*.a": {&rule.MatchingKindRule{Instance: ""}},
	}

	require.NoError(t, err)
	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceComplexMap(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	err = validator.UseReference(reflect.TypeOf(struct {
		C map[string]struct {
			A string `bson:"a"`
		} `bson:"c"`
	}{}))
	expectedRule := map[operation.Path][]rule.Rule{
		"c": {&rule.MatchingKindRule{Instance: map[string]struct {
			A string `bson:"a"`
		}{}}},
	}
	expectedWildcardRules := map[operation.Path][]rule.Rule{
		"c.*": {&rule.MatchingKindRule{Instance: struct {
			A string `bson:"a"`
		}{}}},
		"c.*.a": {&rule.MatchingKindRule{Instance: ""}},
	}

	require.NoError(t, err)
	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceComplexNested(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	err = validator.UseReference(reflect.TypeOf(struct {
		D map[string][]struct {
			A string `bson:"a"`
		} `bson:"d"`
		E []map[string]struct {
			A string `bson:"a"`
		} `bson:"e"`
	}{}))
	expectedRule := map[operation.Path][]rule.Rule{
		"d": {&rule.MatchingKindRule{Instance: map[string][]struct {
			A string `bson:"a"`
		}{}}},
		"e": {&rule.MatchingKindRule{Instance: []map[string]struct {
			A string `bson:"a"`
		}{}}},
	}
	expectedWildcardRules := map[operation.Path][]rule.Rule{
		"d.*": {&rule.MatchingKindRule{Instance: []struct {
			A string `bson:"a"`
		}{}}},
		"d.*.*": {&rule.MatchingKindRule{Instance: struct {
			A string `bson:"a"`
		}{}}},
		"d.*.*.a": {&rule.MatchingKindRule{Instance: ""}},
		"e.*": {&rule.MatchingKindRule{Instance: map[string]struct {
			A string `bson:"a"`
		}{}}},
		"e.*.*": {&rule.MatchingKindRule{Instance: struct {
			A string `bson:"a"`
		}{}}},
		"e.*.*.a": {&rule.MatchingKindRule{Instance: ""}},
	}

	require.NoError(t, err)
	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceWithSimpleRules(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	err = validator.UseReference(reflect.TypeOf(struct {
		A string `bson:"a"`
		B string `bson:"b" jp_disallow:"true"`
		C string `bson:"c" jp_min:"3"`
		D string `bson:"d" jp_max:"3"`
		E string `bson:"e" jp_expression:"^\\w+$"`
		F string `bson:"f" jp_op_allowed:"add,remove"`
		G string `bson:"g" jp_op_disallowed:"add,remove"`
	}{}))
	expectedRule := map[operation.Path][]rule.Rule{
		"a": {&rule.MatchingKindRule{Instance: ""}},
		"b": {&rule.MatchingKindRule{Instance: ""}, &rule.DisallowRule{Disallow: true}},
		"c": {&rule.MatchingKindRule{Instance: ""}, &rule.MinRule{Min: 3}},
		"d": {&rule.MatchingKindRule{Instance: ""}, &rule.MaxRule{Max: 3}},
		"e": {
			&rule.MatchingKindRule{Instance: ""},
			&rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
		},
		"f": {
			&rule.MatchingKindRule{Instance: ""},
			&rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation, operation.RemoveOperation}},
		},
		"g": {
			&rule.MatchingKindRule{Instance: ""},
			&rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.AddOperation, operation.RemoveOperation}},
		},
	}
	expectedWildcardRules := map[operation.Path][]rule.Rule{}

	require.NoError(t, err)
	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceWithHeredityStruct(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	type demoStruct struct {
		DoInherit   string `bson:"do_inherit"`
		DoOverwrite string `bson:"do_overwrite" jp_disallow:"false" jp_min:"2" jp_max:"2" jp_expression:"^\\d+$" jp_op_allowed:"move" jp_op_disallowed:"replace"` //nolint:lll
		Nested      struct {
			DoInherit string `bson:"do_inherit"`
		} `bson:"nested"`
	}
	// TODO Add another level

	err = validator.UseReference(reflect.TypeOf(struct {
		A demoStruct `bson:"a" jp_inherit:"jp_disallow,jp_min,jp_max,jp_expression,jp_op_allowed,jp_op_disallowed" jp_disallow:"true" jp_min:"3" jp_max:"3" jp_expression:"^\\w+$" jp_op_allowed:"add" jp_op_disallowed:"remove"` //nolint:lll
	}{}))
	require.NoError(t, err)

	expectedRule := map[operation.Path][]rule.Rule{
		"a": {
			&rule.MatchingKindRule{Instance: demoStruct{}},
			&rule.DisallowRule{Disallow: true},
			&rule.MinRule{Min: 3}, &rule.MaxRule{Max: 3},
			&rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			&rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			&rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.do_inherit": {
			&rule.MatchingKindRule{Instance: ""},
			&rule.DisallowRule{Disallow: true},
			&rule.MinRule{Min: 3}, &rule.MaxRule{Max: 3},
			&rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			&rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			&rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.do_overwrite": {
			&rule.MatchingKindRule{Instance: ""},
			&rule.DisallowRule{Disallow: false},
			&rule.MinRule{Min: 2}, &rule.MaxRule{Max: 2},
			&rule.ExpressionRule{Expression: `^\d+$`, Regex: *regexp.MustCompile(`^\d+$`)},
			&rule.AllowedOperationsRule{Operations: []operation.Operation{operation.MoveOperation}},
			&rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.ReplaceOperation}},
		},
		"a.nested.do_inherit": {
			&rule.MatchingKindRule{Instance: ""},
			&rule.DisallowRule{Disallow: true},
			&rule.MinRule{Min: 3}, &rule.MaxRule{Max: 3},
			&rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			&rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			&rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
	}
	expectedWildcardRules := map[operation.Path][]rule.Rule{}

	compareMapWithRandomArray(t, expectedRule, validator.rules)
	compareMapWithRandomArray(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceWithHeredityArray(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	type demoStruct struct {
		DoInherit   string `bson:"do_inherit"`
		DoOverwrite string `bson:"do_overwrite" jp_disallow:"false" jp_min:"2" jp_max:"2" jp_expression:"^\\d+$" jp_op_allowed:"move" jp_op_disallowed:"replace"` //nolint:lll
	}

	err = validator.UseReference(reflect.TypeOf(struct {
		A []demoStruct `bson:"a" jp_inherit:"jp_disallow,jp_min,jp_max,jp_expression,jp_op_allowed,jp_op_disallowed" jp_disallow:"true" jp_min:"3" jp_max:"3" jp_expression:"^\\w+$" jp_op_allowed:"add" jp_op_disallowed:"remove"` //nolint:lll
	}{}))
	require.NoError(t, err)

	expectedRule := map[operation.Path][]rule.Rule{
		"a": {
			&rule.MatchingKindRule{Instance: []demoStruct{}},
			&rule.DisallowRule{Disallow: true},
			&rule.MinRule{Min: 3}, &rule.MaxRule{Max: 3},
			&rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			&rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			&rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
	}
	expectedWildcardRules := map[operation.Path][]rule.Rule{
		"a.*": {
			&rule.MatchingKindRule{Instance: demoStruct{}},
			&rule.DisallowRule{Disallow: true},
			&rule.MinRule{Min: 3}, &rule.MaxRule{Max: 3},
			&rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			&rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			&rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.*.do_inherit": {
			&rule.MatchingKindRule{Instance: ""},
			&rule.DisallowRule{Disallow: true},
			&rule.MinRule{Min: 3}, &rule.MaxRule{Max: 3},
			&rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			&rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			&rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.*.do_overwrite": {
			&rule.MatchingKindRule{Instance: ""},
			&rule.DisallowRule{Disallow: false},
			&rule.MinRule{Min: 2},
			&rule.MaxRule{Max: 2},
			&rule.ExpressionRule{Expression: `^\d+$`, Regex: *regexp.MustCompile(`^\d+$`)},
			&rule.AllowedOperationsRule{Operations: []operation.Operation{operation.MoveOperation}},
			&rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.ReplaceOperation}},
		},
	}

	compareMapWithRandomArray(t, expectedRule, validator.rules)
	compareMapWithRandomArray(t, expectedWildcardRules, validator.wildcardRules)
}

// TODO same for  map

func TestUseReferenceWithValidate(t *testing.T) {
	t.Parallel()

	// TODO
}
