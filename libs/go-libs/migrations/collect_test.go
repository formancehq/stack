package migrations

import (
	"fmt"
	"io/fs"
	"math/rand"
	"path/filepath"
	"testing"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCollect(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	migrationsFS := NewMockMigrationFileSystem(ctrl)

	const numberOfFiles = 100

	sortedFiles := make([]string, 0)
	for i := 0; i < numberOfFiles; i++ {
		filename := fmt.Sprintf("%d-migrate.sql", i)
		sortedFiles = append(sortedFiles, filename)

		migrationsFS.EXPECT().ReadFile(filepath.Join("migrations", filename)).Return([]byte(""), nil)
	}

	migrationsFS.EXPECT().
		ReadDir("migrations").
		Return(collectionutils.Map(sortedFiles, func(from string) fs.DirEntry {
			return mockDirEntry(from)
		}), nil)

	// shuffle migration names to ensure the collector sort them
	shuffledFiles := make([]string, numberOfFiles)
	copy(shuffledFiles, sortedFiles)
	rand.Shuffle(len(sortedFiles), func(i, j int) {
		shuffledFiles[i], shuffledFiles[j] = shuffledFiles[j], shuffledFiles[i]
	})

	migrations, err := CollectMigrationFiles(migrationsFS, "migrations")
	require.NoError(t, err)

	require.Equal(t, sortedFiles, collectionutils.Map(migrations, func(from Migration) string {
		return from.Name
	}))
}

type mockDirEntry string

func (m mockDirEntry) Name() string {
	return string(m)
}

func (m mockDirEntry) IsDir() bool {
	return false
}

func (m mockDirEntry) Type() fs.FileMode {
	return 0
}

func (m mockDirEntry) Info() (fs.FileInfo, error) {
	return nil, nil
}

var _ fs.DirEntry = (*mockDirEntry)(nil)
