package rule

import (
	"reflect"
	"testing"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
	"github.com/stretchr/testify/require"
)

func TestRuleDisallowedOperation(t *testing.T) {
	t.Parallel()

	rule := DisallowedOperationsRule{}
	err := rule.UseValue("a", reflect.Array, nil, "replace,move")

	require.NoError(t, err)

	t.Run("AllowedCase", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, rule.Apply(operation.Spec{
			Operation: operation.AddOperation,
		}))
	})

	t.Run("ForbiddenCase", func(t *testing.T) {
		t.Parallel()
		require.Error(t, rule.Apply(operation.Spec{
			Operation: operation.ReplaceOperation,
		}))
	})
}
