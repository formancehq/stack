package v2

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/formancehq/payments/internal/api/backend"
	"github.com/formancehq/payments/internal/otel"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
)

func connectorsInstall(backend backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.Tracer().Start(r.Context(), "v2_connectorsInstall")
		defer span.End()

		provider := strings.ToLower(connectorProvider(r))
		if provider == "" {
			otel.RecordError(span, errors.New("provider is required"))
			api.BadRequest(w, ErrValidation, errors.New("provider is required"))
			return
		}

		config, err := io.ReadAll(r.Body)
		if err != nil {
			otel.RecordError(span, err)
			api.BadRequest(w, ErrMissingOrInvalidBody, err)
			return
		}

		// Detach the context to avoid cancellation of the installation process
		// leading to a partial installation
		ctx, _ = contextutil.Detached(ctx)
		connectorID, err := backend.ConnectorsInstall(ctx, provider, config)
		if err != nil {
			otel.RecordError(span, err)
			handleServiceErrors(w, r, err)
			return
		}

		api.Created(w, connectorID.String())
	}
}
