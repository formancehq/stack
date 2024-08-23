package runner

import (
	"context"

	ingester "github.com/formancehq/stack/ee/ingester/internal"
)

//go:generate mockgen -source store.go -destination store_generated.go -package runner . Store
type Store interface {
	ListEnabledPipelines(ctx context.Context) ([]ingester.Pipeline, error)
	StoreState(ctx context.Context, id string, state ingester.State) error
}
