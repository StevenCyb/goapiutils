package jsonpatch

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidPath(t *testing.T) {
	t.Parallel()

	t.Run("ShortPath_Success", func(t *testing.T) {
		t.Parallel()

		path := Path("root")
		require.True(t, path.Valid())
	})

	t.Run("LongPath_Success", func(t *testing.T) {
		t.Parallel()

		path := Path("parent.child")
		require.True(t, path.Valid())
	})

	t.Run("InvalidPath_Fail", func(t *testing.T) {
		t.Parallel()

		path := Path("now..valid")
		require.False(t, path.Valid())

		path = Path(".fds.")
		require.False(t, path.Valid())

		path = Path("")
		require.False(t, path.Valid())

		path = Path(".##")
		require.False(t, path.Valid())
	})
}
