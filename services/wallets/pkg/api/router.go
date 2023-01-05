package api

import (
	"net/http"

	sharedhealth "github.com/formancehq/go-libs/health"
	"github.com/formancehq/wallets/pkg/wallet"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
)

func NewRouter(
	funding *wallet.FundingService,
	repository *wallet.Repository,
	healthController *sharedhealth.HealthController,
) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/_healthcheck", healthController.Check)
	r.Group(func(r chi.Router) {
		r.Use(otelchi.Middleware("wallets"))
		r.Use(middleware.Logger)
		r.Use(middleware.AllowContentType("application/json"))
		r.Use(func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				handler.ServeHTTP(w, r)
			})
		})
		main := NewMainHandler(funding, repository)

		r.Route("/wallets", func(r chi.Router) {
			r.Get("/", main.ListWalletsHandler)
			r.Post("/", main.CreateWalletHandler)
			r.Route("/{walletID}", func(r chi.Router) {
				r.Get("/", main.GetWalletHandler)
				r.Patch("/", main.PatchWalletHandler)
				r.Post("/debit", main.DebitWalletHandler)
				r.Post("/credit", main.CreditWalletHandler)
			})
		})
		r.Route("/transactions", func(r chi.Router) {
			r.Get("/", main.ListTransactions)
		})
		r.Route("/holds", func(r chi.Router) {
			r.Get("/", main.ListHoldsHandler)
			r.Route("/{holdID}", func(r chi.Router) {
				r.Get("/", main.GetHoldHandler)
				r.Post("/confirm", main.ConfirmHoldHandler)
				r.Post("/void", main.VoidHoldHandler)
			})
		})
	})

	return r
}
