package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadataRegexp(t *testing.T) {
	t.Parallel()

	t.Run("valid tests", func(t *testing.T) {
		t.Parallel()
		require.True(t, metadataRegex.MatchString("metadata[foo]"))
		require.True(t, metadataRegex.MatchString("metadata[foo_bar]"))
		require.True(t, metadataRegex.MatchString("metadata[foo/bar]"))
		require.True(t, metadataRegex.MatchString("metadata[foo.bar]"))
	})

	t.Run("invalid tests", func(t *testing.T) {
		t.Parallel()

		require.False(t, metadataRegex.MatchString("metadata[foo"))
		require.False(t, metadataRegex.MatchString("metadata/foo"))
		require.False(t, metadataRegex.MatchString("metadata.foo"))
	})
}
