package migrations

import (
	"context"
	"embed"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func CollectMigrationFiles(fs embed.FS) ([]Migration, error) {
	entries, err := fs.ReadDir("migrations")
	if err != nil {
		return nil, errors.Wrap(err, "collecting migration files")
	}

	ret := make([]Migration, len(entries))
	for i, entry := range entries {
		fileContent, err := fs.ReadFile(filepath.Join("migrations", entry.Name()))
		if err != nil {
			return nil, errors.Wrapf(err, "reading migration file %s", entry.Name())
		}

		ret[i] = Migration{
			Name: entry.Name(),
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.ExecContext(ctx, string(fileContent))
				return err
			},
		}
	}

	return ret, nil
}
