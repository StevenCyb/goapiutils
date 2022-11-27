package rule

import (
	"testing"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
	"github.com/stretchr/testify/require"
)

type objectA struct {
	Mapping map[string]struct {
		D string `bson:"d"`
	} `bson:"mapping"`
	Nested struct {
		B *string `bson:"b"`
		A string  `bson:"a"`
	} `bson:"nested"`
	Name   string `bson:"name"`
	IntArr []int  `bson:"int_arr"`
	ObjArr []struct {
		C string `bson:"c"`
	} `bson:"obj_arr"`
	Age int `bson:"age"`
}

//nolint:unused,revive,stylecheck
type objectB struct {
	mapping map[string]struct{ d string }
	nested  struct {
		a string
		b string
	}
	name    string
	int_arr []int
	obj_arr []struct{ c string }
	age     int
}

func TestRuleMatchingKindEqualType(t *testing.T) {
	t.Parallel()

	rule := MatchingKindRule{instance: ""}
	require.NoError(t, rule.Apply(operation.Spec{Value: "hello"}))

	rule = MatchingKindRule{instance: uint32(0)}
	require.NoError(t, rule.Apply(operation.Spec{Value: uint32(4)}))

	rule = MatchingKindRule{instance: []int{}}
	require.NoError(t, rule.Apply(operation.Spec{Value: []int{1, 2, 3}}))

	rule = MatchingKindRule{instance: objectA{}}
	require.NoError(t, rule.Apply(operation.Spec{Value: objectB{}}))
}

func TestRuleMatchingKindNotEqualType(t *testing.T) {
	t.Parallel()

	rule := MatchingKindRule{instance: ""}
	err := rule.Apply(operation.Spec{Value: 1})
	require.Error(t, err)
	require.Equal(t, "'root value' has invalid kind 'int', must be 'string'", err.Error())

	rule = MatchingKindRule{instance: []string{}}
	err = rule.Apply(operation.Spec{Value: []int{1, 2, 3}})
	require.Error(t, err)
	require.Equal(t, "'root value item' has invalid kind 'int', must be 'string'", err.Error())

	rule = MatchingKindRule{instance: objectA{}}
	err = rule.Apply(operation.Spec{Value: struct {
		aa string
	}{}})
	require.Error(t, err)
	require.Equal(t, "unknown field 'aa'", err.Error())

	rule = MatchingKindRule{instance: objectA{}}
	err = rule.Apply(operation.Spec{Value: struct {
		nested struct {
			c string
		}
	}{}})
	require.Error(t, err)
	require.Equal(t, "unknown field 'c'", err.Error())

	rule = MatchingKindRule{instance: objectA{}}
	err = rule.Apply(operation.Spec{Value: struct {
		obj_arr []struct { //nolint:revive,stylecheck
			e string
		}
	}{}})
	require.Error(t, err)
	require.Equal(t, "unknown field 'e'", err.Error())

	rule = MatchingKindRule{instance: objectA{}}
	err = rule.Apply(operation.Spec{Value: struct {
		mapping map[string]struct {
			e string
		}
	}{}})
	require.Error(t, err)
	require.Equal(t, "unknown field 'e'", err.Error())
}
