package drivers

import (
	"context"

	ingester "github.com/formancehq/stack/ee/ingester/internal"
)

//go:generate mockgen -source driver.go -destination driver_generated.go -package drivers . Driver
type Driver interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Accept(ctx context.Context, logs ...ingester.LogWithModule) ([]error, error)
}
