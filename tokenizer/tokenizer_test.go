package tokenizer

import (
	"fmt"
	"testing"

	"github.com/StevenCyb/goquery/errs"

	"github.com/stretchr/testify/require"
)

func TestTokenizer(t *testing.T) {
	var TYPE_NONE Type = "NONE"
	var TYPE_EQUAL Type = "EQUAL"
	var TYPE_WORD Type = "WORD"
	key := "hello"
	separator := "="
	value := "world"

	t.Run("SpecParsing", func(t *testing.T) {
		tokenizer := NewTokenizer(
			fmt.Sprintf("%s%s%s", key, separator, value),
			TYPE_NONE, TYPE_NONE,
			[]*Spec{
				NewSpec("^=", TYPE_EQUAL),
				NewSpec("^[a-z]+", TYPE_WORD),
			},
			nil)

		token, err := tokenizer.GetNextToken()
		require.NoError(t, err)
		require.Equal(t, TYPE_WORD, token.Type)
		require.Equal(t, key, token.Value)

		token, err = tokenizer.GetNextToken()
		require.NoError(t, err)
		require.Equal(t, TYPE_EQUAL, token.Type)
		require.Equal(t, separator, token.Value)

		token, err = tokenizer.GetNextToken()
		require.NoError(t, err)
		require.Equal(t, TYPE_WORD, token.Type)
		require.Equal(t, value, token.Value)
	})

	t.Run("SpecParsingWithSkip", func(t *testing.T) {
		var TYPE_SKIP Type = "SKIP"
		tokenizer := NewTokenizer(
			`  hello  = world `,
			TYPE_SKIP, TYPE_NONE,
			[]*Spec{
				NewSpec(`^\s+`, TYPE_SKIP),
				NewSpec("^=", TYPE_EQUAL),
				NewSpec("^[a-z]+", TYPE_WORD),
			},
			nil)

		token, err := tokenizer.GetNextToken()
		require.NoError(t, err)
		require.Equal(t, TYPE_WORD, token.Type)
		require.Equal(t, key, token.Value)

		token, err = tokenizer.GetNextToken()
		require.NoError(t, err)
		require.Equal(t, TYPE_EQUAL, token.Type)
		require.Equal(t, separator, token.Value)

		token, err = tokenizer.GetNextToken()
		require.NoError(t, err)
		require.Equal(t, TYPE_WORD, token.Type)
		require.Equal(t, value, token.Value)
	})

	t.Run("PolicyWithoutViolation", func(t *testing.T) {
		t.Run("SpecParsing", func(t *testing.T) {
			tokenizer := NewTokenizer(
				fmt.Sprintf("%s%s%s", key, separator, value),
				TYPE_NONE, TYPE_WORD,
				[]*Spec{
					NewSpec("^=", TYPE_EQUAL),
					NewSpec("^[a-z]+", TYPE_WORD),
				},
				NewPolicy(WHITELIST_POLICY, key, value))

			token, err := tokenizer.GetNextToken()
			require.NoError(t, err)
			require.Equal(t, TYPE_WORD, token.Type)
			require.Equal(t, key, token.Value)

			token, err = tokenizer.GetNextToken()
			require.NoError(t, err)
			require.Equal(t, TYPE_EQUAL, token.Type)
			require.Equal(t, separator, token.Value)

			token, err = tokenizer.GetNextToken()
			require.NoError(t, err)
			require.Equal(t, TYPE_WORD, token.Type)
			require.Equal(t, value, token.Value)
		})
	})

	t.Run("PolicyWithViolation", func(t *testing.T) {
		t.Run("SpecParsing", func(t *testing.T) {
			tokenizer := NewTokenizer(
				fmt.Sprintf("%s%s%s", key, separator, value),
				TYPE_NONE, TYPE_WORD,
				[]*Spec{
					NewSpec("^=", TYPE_EQUAL),
					NewSpec("^[a-z]+", TYPE_WORD),
				},
				NewPolicy(WHITELIST_POLICY, key))

			token, err := tokenizer.GetNextToken()
			require.NoError(t, err)
			require.Equal(t, TYPE_WORD, token.Type)
			require.Equal(t, key, token.Value)

			token, err = tokenizer.GetNextToken()
			require.NoError(t, err)
			require.Equal(t, TYPE_EQUAL, token.Type)
			require.Equal(t, separator, token.Value)

			_, err = tokenizer.GetNextToken()
			require.Equal(t, errs.NewErrPolicyViolation(value), err)
		})
	})
}
