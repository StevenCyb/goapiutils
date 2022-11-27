package rule

import (
	"reflect"
	"testing"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
	"github.com/stretchr/testify/require"
)

func TestRuleDisallow(t *testing.T) {
	t.Parallel()

	t.Run("AllowedCase", func(t *testing.T) {
		t.Parallel()

		rule := DisallowRule{}
		err := rule.UseValue("a", reflect.Bool, nil, "false")

		require.NoError(t, err)
		require.NoError(t, rule.Apply(operation.Spec{
			Path: "a.c",
		}))
	})

	t.Run("ForbiddenCase", func(t *testing.T) {
		t.Parallel()

		rule := DisallowRule{}
		err := rule.UseValue("a", reflect.Bool, nil, "true")

		require.NoError(t, err)
		require.Error(t, rule.Apply(operation.Spec{
			Path: "a.b",
		}))
	})
}
