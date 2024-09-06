package drivers

import (
	"context"

	ingester "github.com/formancehq/stack/ee/ingester/internal"
)

//go:generate mockgen -source store.go -destination store_generated.go -package drivers . Store
type Store interface {
	GetConnector(ctx context.Context, id string) (*ingester.Connector, error)
}
