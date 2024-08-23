package api

import (
	"testing"

	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/logging"
)

func newAPI(t *testing.T, backend Backend) *API {
	t.Helper()

	return NewAPI(
		backend,
		health.NewHealthController(),
		auth.NewNoAuth(),
		logging.Testing(),
		api.ServiceInfo{
			Version: "testing",
			Debug:   testing.Verbose(),
		},
	)
}
