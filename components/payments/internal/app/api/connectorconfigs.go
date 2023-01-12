package api

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/payments/internal/app/connectors/bankingcircle"
	"github.com/formancehq/payments/internal/app/connectors/currencycloud"
	"github.com/formancehq/payments/internal/app/connectors/dummypay"
	"github.com/formancehq/payments/internal/app/connectors/modulr"
	"github.com/formancehq/payments/internal/app/connectors/stripe"
	"github.com/formancehq/payments/internal/app/connectors/wise"

	"github.com/formancehq/payments/internal/pkg/configtemplate"
)

func connectorConfigsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: It's not ideal to re-identify available connectors
		// Refactor it when refactoring the HTTP lib.

		configs := configtemplate.BuildConfigs(
			bankingcircle.Config{},
			currencycloud.Config{},
			dummypay.Config{},
			modulr.Config{},
			stripe.Config{},
			wise.Config{},
		)

		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(configs)
		if err != nil {
			handleServerError(w, r, err)

			return
		}
	}
}
