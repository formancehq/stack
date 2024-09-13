package engine

import (
	"github.com/formancehq/payments/internal/connectors/engine/plugins"
	"github.com/pkg/errors"
)

var (
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("not found")
)

func handlePluginError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, plugins.ErrNotFound):
		return errors.Wrap(ErrNotFound, err.Error())
	default:
		return err
	}
}
