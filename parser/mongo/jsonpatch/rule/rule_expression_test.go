package rule

import (
	"reflect"
	"testing"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
	"github.com/stretchr/testify/require"
)

func TestRuleExpression(t *testing.T) {
	t.Parallel()

	rule := ExpressionRule{}
	err := rule.UseValue("a", reflect.Array, nil, "^[a-z]+$")

	require.NoError(t, err)

	t.Run("AllowedCase", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, rule.Apply(operation.Spec{
			Value: "hello",
		}))
	})

	t.Run("ForbiddenCase", func(t *testing.T) {
		t.Parallel()
		require.Error(t, rule.Apply(operation.Spec{
			Value: "123",
		}))
	})
}
