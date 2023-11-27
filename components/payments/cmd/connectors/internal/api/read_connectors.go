package api

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/payments/cmd/connectors/internal/api/backend"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
)

type readConnectorsResponseElement struct {
	Provider    models.ConnectorProvider `json:"provider" bson:"provider"`
	ConnectorID string                   `json:"connectorID" bson:"connectorID"`
	Name        string                   `json:"name" bson:"name"`
}

func readConnectorsHandler(b backend.ServiceBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := b.GetService().ListConnectors(r.Context())
		if err != nil {
			handleServiceErrors(w, r, err)
			return
		}

		data := make([]readConnectorsResponseElement, len(res))

		for i := range res {
			data[i] = readConnectorsResponseElement{
				Provider:    res[i].Provider,
				ConnectorID: res[i].ID.String(),
				Name:        res[i].Name,
			}
		}

		err = json.NewEncoder(w).Encode(
			api.BaseResponse[[]readConnectorsResponseElement]{
				Data: &data,
			})
		if err != nil {
			panic(err)
		}
	}
}
