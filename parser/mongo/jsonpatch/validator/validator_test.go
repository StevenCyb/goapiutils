//nolint:dupl
package validator

import (
	"errors"
	"reflect"
	"regexp"
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

		err := validator.RegisterRule("jp_name", demoRule{})
		require.NoError(t, err)
	})

	t.Run("DuplicateRuleName_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.RegisterRule("jp_disallow", demoRule{})
		require.Error(t, err)
		require.Equal(t, ErrDuplicateRuleTags, err)
	})

	t.Run("MissingPrefix_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.RegisterRule("name", demoRule{})
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
	expectedRule := map[operation.Path]map[string]rule.Rule{
		"b": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""}},
		"c": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: 0}},
		"d": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: []string{}}},
		"e": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: []int{}}},
		"f": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: map[string]string{}}},
		"g": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: map[string]int{}}},
		"h": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: struct {
			A string
			B string `bson:"b"`
		}{}}},
		"h.b": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""}},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{
		"d.*": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""}},
		"e.*": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: 0}},
		"f.*": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""}},
		"g.*": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: 0}},
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
	expectedRule := map[operation.Path]map[string]rule.Rule{
		"a": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: struct {
			A struct {
				B string `bson:"b"`
			} `bson:"a"`
		}{}}},
		"a.a": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: struct {
			B string `bson:"b"`
		}{}}},
		"a.a.b": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""}},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{}

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
	expectedRule := map[operation.Path]map[string]rule.Rule{
		"b": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: []struct {
			A string `bson:"a"`
		}{}}},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{
		"b.*": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: struct {
			A string `bson:"a"`
		}{}}},
		"b.*.a": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""}},
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
	expectedRule := map[operation.Path]map[string]rule.Rule{
		"c": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: map[string]struct {
			A string `bson:"a"`
		}{}}},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{
		"c.*": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: struct {
			A string `bson:"a"`
		}{}}},
		"c.*.a": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""}},
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
	expectedRule := map[operation.Path]map[string]rule.Rule{
		"d": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: map[string][]struct {
			A string `bson:"a"`
		}{}}},
		"e": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: []map[string]struct {
			A string `bson:"a"`
		}{}}},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{
		"d.*": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: []struct {
			A string `bson:"a"`
		}{}}},
		"d.*.*": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: struct {
			A string `bson:"a"`
		}{}}},
		"d.*.*.a": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""}},
		"e.*": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: map[string]struct {
			A string `bson:"a"`
		}{}}},
		"e.*.*": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: struct {
			A string `bson:"a"`
		}{}}},
		"e.*.*.a": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""}},
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
	expectedRule := map[operation.Path]map[string]rule.Rule{
		"a": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""}},
		"b": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""},
			"jp_disallow":              &rule.DisallowRule{Disallow: true},
		},
		"c": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""}, "jp_min": &rule.MinRule{Min: 3}},
		"d": {"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""}, "jp_max": &rule.MaxRule{Max: 3}},
		"e": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""},
			"jp_expression":            &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
		},
		"f": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""},
			"jp_op_allowed": &rule.AllowedOperationsRule{
				Operations: []operation.Operation{operation.AddOperation, operation.RemoveOperation},
			},
		},
		"g": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{
				Operations: []operation.Operation{operation.AddOperation, operation.RemoveOperation},
			},
		},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{}

	require.NoError(t, err)
	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceWithHeredityStruct(t *testing.T) { //nolint:funlen
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	type demoStruct struct {
		DoInherit   string `bson:"do_inherit"`
		DoOverwrite string `bson:"do_overwrite" jp_disallow:"false" jp_min:"2" jp_max:"2" jp_expression:"^\\d+$" jp_op_allowed:"move" jp_op_disallowed:"replace"` //nolint:lll
		Nested      struct {
			DoInherit string `bson:"do_inherit"`
			Nested    struct {
				DoInherit string `bson:"do_inherit"`
			} `bson:"nested"`
		} `bson:"nested"`
	}

	err = validator.UseReference(reflect.TypeOf(struct {
		A demoStruct `bson:"a" jp_inherit:"jp_disallow,jp_min,jp_max,jp_expression,jp_op_allowed,jp_op_disallowed" jp_disallow:"true" jp_min:"3" jp_max:"3" jp_expression:"^\\w+$" jp_op_allowed:"add" jp_op_disallowed:"remove"` //nolint:lll
	}{}))
	require.NoError(t, err)

	expectedRule := map[operation.Path]map[string]rule.Rule{
		"a": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: demoStruct{}},
			"jp_disallow":              &rule.DisallowRule{Disallow: true},
			"jp_min":                   &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.do_inherit": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""},
			"jp_disallow":              &rule.DisallowRule{Disallow: true},
			"jp_min":                   &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.do_overwrite": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""},
			"jp_disallow":              &rule.DisallowRule{Disallow: false},
			"jp_min":                   &rule.MinRule{Min: 2}, "jp_max": &rule.MaxRule{Max: 2},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\d+$`, Regex: *regexp.MustCompile(`^\d+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.MoveOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.ReplaceOperation}},
		},
		"a.nested": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: struct {
				DoInherit string `bson:"do_inherit"`
				Nested    struct {
					DoInherit string `bson:"do_inherit"`
				} `bson:"nested"`
			}{}},
			"jp_disallow": &rule.DisallowRule{Disallow: true},
			"jp_min":      &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.nested.do_inherit": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""},
			"jp_disallow":              &rule.DisallowRule{Disallow: true},
			"jp_min":                   &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.nested.nested": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: struct {
				DoInherit string `bson:"do_inherit"`
			}{}},
			"jp_disallow": &rule.DisallowRule{Disallow: true},
			"jp_min":      &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.nested.nested.do_inherit": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""},
			"jp_disallow":              &rule.DisallowRule{Disallow: true},
			"jp_min":                   &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{}

	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
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

	expectedRule := map[operation.Path]map[string]rule.Rule{
		"a": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: []demoStruct{}},
			"jp_disallow":              &rule.DisallowRule{Disallow: true},
			"jp_min":                   &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{
		"a.*": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: demoStruct{}},
			"jp_disallow":              &rule.DisallowRule{Disallow: true},
			"jp_min":                   &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.*.do_inherit": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""},
			"jp_disallow":              &rule.DisallowRule{Disallow: true},
			"jp_min":                   &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.*.do_overwrite": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""},
			"jp_disallow":              &rule.DisallowRule{Disallow: false},
			"jp_min":                   &rule.MinRule{Min: 2}, "jp_max": &rule.MaxRule{Max: 2},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\d+$`, Regex: *regexp.MustCompile(`^\d+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.MoveOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.ReplaceOperation}},
		},
	}

	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestUseReferenceWithHeredityMap(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct{}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	type demoStruct struct {
		DoInherit   string `bson:"do_inherit"`
		DoOverwrite string `bson:"do_overwrite" jp_disallow:"false" jp_min:"2" jp_max:"2" jp_expression:"^\\d+$" jp_op_allowed:"move" jp_op_disallowed:"replace"` //nolint:lll
	}

	err = validator.UseReference(reflect.TypeOf(struct {
		A map[string]demoStruct `bson:"a" jp_inherit:"jp_disallow,jp_min,jp_max,jp_expression,jp_op_allowed,jp_op_disallowed" jp_disallow:"true" jp_min:"3" jp_max:"3" jp_expression:"^\\w+$" jp_op_allowed:"add" jp_op_disallowed:"remove"` //nolint:lll
	}{}))
	require.NoError(t, err)

	expectedRule := map[operation.Path]map[string]rule.Rule{
		"a": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: map[string]demoStruct{}},
			"jp_disallow":              &rule.DisallowRule{Disallow: true},
			"jp_min":                   &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
	}
	expectedWildcardRules := map[operation.Path]map[string]rule.Rule{
		"a.*": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: demoStruct{}},
			"jp_disallow":              &rule.DisallowRule{Disallow: true},
			"jp_min":                   &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.*.do_inherit": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""},
			"jp_disallow":              &rule.DisallowRule{Disallow: true},
			"jp_min":                   &rule.MinRule{Min: 3}, "jp_max": &rule.MaxRule{Max: 3},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\w+$`, Regex: *regexp.MustCompile(`^\w+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.AddOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.RemoveOperation}},
		},
		"a.*.do_overwrite": {
			"jp_general_matching_kind": &rule.MatchingKindRule{Instance: ""},
			"jp_disallow":              &rule.DisallowRule{Disallow: false},
			"jp_min":                   &rule.MinRule{Min: 2}, "jp_max": &rule.MaxRule{Max: 2},
			"jp_expression":    &rule.ExpressionRule{Expression: `^\d+$`, Regex: *regexp.MustCompile(`^\d+$`)},
			"jp_op_allowed":    &rule.AllowedOperationsRule{Operations: []operation.Operation{operation.MoveOperation}},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{Operations: []operation.Operation{operation.ReplaceOperation}},
		},
	}

	require.Equal(t, expectedRule, validator.rules)
	require.Equal(t, expectedWildcardRules, validator.wildcardRules)
}

func TestValidateTypecheckRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct {
		A string `bson:"a"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidType_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "a", Value: "new"})
		require.NoError(t, err)
	})

	t.Run("InvalidType_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "a", Value: 123})
		require.Error(t, err)
		require.Equal(t, "operation no allowed: 'root value' has invalid kind 'int', must be 'string'", err.Error())
	})
}

type demoValidateStruct struct { //nolint:govet
	A int `bson:"a"`
	B struct {
		C int `bson:"c"`
	} `bson:"b"`
	D []struct {
		E int `bson:"e"`
	} `bson:"d"`
	F map[string]struct {
		G int `bson:"g"`
	} `bson:"f"`
}

func TestValidatePath(t *testing.T) { //nolint:funlen
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(demoValidateStruct{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidRootPath_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "a", Value: 1})
		require.NoError(t, err)
	})

	t.Run("ValidNestedPath_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "b.c", Value: 1})
		require.NoError(t, err)
	})

	t.Run("ValidNestedArrayPath_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "d.0.e", Value: 1})
		require.NoError(t, err)
	})

	t.Run("ValidNestedMapPath_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "f.0.g", Value: 1})
		require.NoError(t, err)
	})

	t.Run("InvalidRootPath_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "aa", Value: 1})
		require.Error(t, err)
		require.Equal(t, "defined path 'aa' is unknown", err.Error())
	})

	t.Run("InvalidNestedPath_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "x.f", Value: 1})
		require.Error(t, err)
		require.Equal(t, "defined path 'x.f' is unknown", err.Error())
	})

	t.Run("InvalidNestedPath_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "d.0.c", Value: 1})
		require.Error(t, err)
		require.Equal(t, "defined path 'd.0.c' is unknown", err.Error())
	})

	t.Run("InvalidNestedPath_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "f.0.c", Value: 1})
		require.Error(t, err)
		require.Equal(t, "defined path 'f.0.c' is unknown", err.Error())
	})
}

func TestValidateDisallowRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct {
		A string `bson:"a" jp_disallow:"false"`
		B string `bson:"b" jp_disallow:"true"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidType_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "a", Value: "new"})
		require.NoError(t, err)
	})

	t.Run("InvalidType_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "b", Value: 123})
		require.Error(t, err)
		require.Equal(t, "operation no allowed: 'root value' has invalid kind 'int', must be 'string'", err.Error())
	})
}

func TestValidateMinRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct {
		A int `bson:"a" jp_min:"3"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidType_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "a", Value: 3})
		require.NoError(t, err)
	})

	t.Run("InvalidType_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "a", Value: 1})
		require.Error(t, err)
		require.Equal(t, "operation no allowed: value is less then specified: '1.000000' < '3.000000'", err.Error())
	})
}

func TestValidateMaxRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct {
		A int `bson:"a" jp_max:"3"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidType_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "a", Value: 3})
		require.NoError(t, err)
	})

	t.Run("InvalidType_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "a", Value: 123})
		require.Error(t, err)
		require.Equal(t, "operation no allowed: value is greater then specified: '123.000000' > '3.000000'", err.Error())
	})
}

func TestValidateExpressionRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct {
		A string `bson:"a" jp_expression:"^[a-z]+$"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidType_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "a", Value: "abc"})
		require.NoError(t, err)
	})

	t.Run("InvalidType_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "a", Value: "123"})
		require.Error(t, err)
		require.Equal(t, "operation no allowed: expression '^[a-z]+$' not match 123", err.Error())
	})
}

func TestValidateAllowedOperationsRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct {
		A string `bson:"a" jp_op_allowed:"add"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidType_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "a", Value: "abc"})
		require.NoError(t, err)
	})

	t.Run("InvalidType_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.RemoveOperation, Path: "a"})
		require.Error(t, err)
		require.Equal(t, "operation no allowed: operation 'remove' not allowed", err.Error())
	})
}

func TestValidateDisallowedOperationsRule(t *testing.T) {
	t.Parallel()

	validator, err := NewValidator(reflect.TypeOf(struct {
		A string `bson:"a" jp_op_disallowed:"remove"`
	}{}))
	require.NoError(t, err)
	require.NotNil(t, validator)

	t.Run("ValidType_Success", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.AddOperation, Path: "a", Value: "abc"})
		require.NoError(t, err)
	})

	t.Run("InvalidType_Fail", func(t *testing.T) {
		t.Parallel()

		err := validator.Validate(operation.Spec{Operation: operation.RemoveOperation, Path: "a"})
		require.Error(t, err)
		require.Equal(t, "operation no allowed: operation 'remove' not allowed", err.Error())
	})
}
