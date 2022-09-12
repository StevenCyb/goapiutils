// nolint:goconst
package patchoperation

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPolicyInterface(t *testing.T) {
	t.Parallel()
	// just test the assigning

	var (
		policy1 Policy = DisallowPathPolicy{}
		policy2 Policy = DisallowOperationOnPathPolicy{}
		policy3 Policy = ForceTypeOnPathPolicy{}
		policy4 Policy = ForceRegexMatchPolicy{}
	)

	// so the variables are used
	require.NotNil(t, policy1)
	require.NotNil(t, policy2)
	require.NotNil(t, policy3)
	require.NotNil(t, policy4)
}

func TestDisallowPathPolicy(t *testing.T) {
	t.Parallel()

	name := "Name"
	forbiddenPath := Path("/user/password")
	allowedPath := Path("/user/username")
	policy := DisallowPathPolicy{Details: name, Path: forbiddenPath}

	require.Equal(t, name, policy.Details)
	require.False(t, policy.Test(OperationSpec{Path: forbiddenPath}))
	require.True(t, policy.Test(OperationSpec{Path: allowedPath}))
}

func TestDisallowOperationOnPathPolicy(t *testing.T) {
	t.Parallel()

	name := "Name"
	path := Path("/user/username")
	notEffected := Path("/user/disabled")
	policy := DisallowOperationOnPathPolicy{
		Details: name, Path: path, Operation: RemoveOperation,
	}

	require.Equal(t, name, policy.Details)
	require.False(t, policy.Test(OperationSpec{Path: path, Operation: RemoveOperation}))
	require.True(t, policy.Test(OperationSpec{Path: path, Operation: ReplaceOperation}))
	require.True(t, policy.Test(OperationSpec{Path: notEffected, Operation: RemoveOperation}))
}

func TestForceTypeOnPathPolicy(t *testing.T) {
	t.Parallel()

	name := "Name"
	path := Path("/product/price")
	notEffected := Path("/product/tag")
	policy := ForceTypeOnPathPolicy{
		Details: name, Path: path, Kind: reflect.Float64,
	}

	require.Equal(t, name, policy.Details)
	require.False(t, policy.Test(OperationSpec{Path: path, Value: "not a price"}))
	require.True(t, policy.Test(OperationSpec{Path: path, Value: 6.99}))
	require.True(t, policy.Test(OperationSpec{Path: notEffected, Value: "vegetable"}))
}

func TestForceRegexMatchPolicy(t *testing.T) {
	t.Parallel()

	name := "Name"
	path := Path("/product/version")
	notEffected := Path("/product/tag")
	policy := ForceRegexMatchPolicy{
		Details: name, Path: path,
		Expression: *regexp.MustCompile(`^v[0-9]*\.[0-9]*\.[0-9]*$`),
	}

	require.Equal(t, name, policy.Details)
	require.False(t, policy.Test(OperationSpec{Path: path, Value: "va.0.0"}))
	require.True(t, policy.Test(OperationSpec{Path: path, Value: "v0.3.7"}))
	require.True(t, policy.Test(OperationSpec{Path: notEffected, Value: "backend"}))
}
