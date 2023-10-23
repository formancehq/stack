package api

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/bankingcircle"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/configtemplate"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/currencycloud"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/dummypay"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/modulr"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/moneycorp"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/stripe"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/wise"
	"github.com/formancehq/stack/libs/go-libs/api"
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
			mangopay.Config{},
			moneycorp.Config{},
		)

		err := json.NewEncoder(w).Encode(api.BaseResponse[configtemplate.Configs]{
			Data: &configs,
		})
		if err != nil {
			handleServerError(w, r, err)

			return
		}
	}
}
