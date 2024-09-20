package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/formancehq/go-libs/health"
	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/stack/ee/stargate/internal/client/controllers"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(
	logger logging.Logger,
	healthController *health.HealthController,
	stargateController *controllers.StargateController,
) chi.Router {
	router := chi.NewMux()

	router.Use(
		cors.New(cors.Options{
			AllowOriginFunc: func(r *http.Request, origin string) bool {
				return true
			},
			AllowCredentials: true,
		}).Handler,
		middleware.Recoverer,
	)

	router.Get("/_healthcheck", healthController.Check)
	router.Get("/_info", stargateController.GetInfo)

	return router
}
