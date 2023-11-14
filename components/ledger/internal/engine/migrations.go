package engine

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/migrations"
)

func (l *Ledger) GetMigrationsInfo(ctx context.Context) ([]migrations.Info, error) {
	return l.bucket.GetMigrationsInfo(ctx)
}
