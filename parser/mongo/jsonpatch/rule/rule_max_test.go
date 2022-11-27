//nolint:dupl
package rule

import (
	"reflect"
	"testing"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
	"github.com/stretchr/testify/require"
)

func TestRuleMax(t *testing.T) {
	t.Parallel()

	rule := MaxRule{}
	err := rule.UseValue("a", reflect.Array, nil, "3")

	require.NoError(t, err)

	t.Run("AllowedCase", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, rule.Apply(operation.Spec{
			Value: 3,
		}))
		require.NoError(t, rule.Apply(operation.Spec{
			Value: uint64(2),
		}))
		require.NoError(t, rule.Apply(operation.Spec{
			Value: 1.2,
		}))
	})

	t.Run("ForbiddenCase", func(t *testing.T) {
		t.Parallel()
		require.Error(t, rule.Apply(operation.Spec{
			Value: 3.4,
		}))
	})
}
