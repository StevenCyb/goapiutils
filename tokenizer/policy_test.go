package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWhitelist(t *testing.T) {
	policy := NewPolicy(WHITELIST_POLICY, "a", "b")
	require.True(t, policy.Allow("a"))
	require.True(t, policy.Allow("b"))
	require.False(t, policy.Allow("c"))
}

func TestBlacklist(t *testing.T) {
	policy := NewPolicy(BLACKLIST_POLICY, "a", "b")
	require.False(t, policy.Allow("a"))
	require.False(t, policy.Allow("b"))
	require.True(t, policy.Allow("c"))
}
