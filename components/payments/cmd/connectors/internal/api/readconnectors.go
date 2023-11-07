package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
)

type readConnectorsRepository interface {
	ListConnectors(ctx context.Context) ([]*models.Connector, error)
}

type readConnectorsResponseElement struct {
	Provider    models.ConnectorProvider `json:"provider" bson:"provider"`
	ConnectorID string                   `json:"connectorID" bson:"connectorID"`
	Name        string                   `json:"name" bson:"name"`
}

func readConnectorsHandler(repo readConnectorsRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := repo.ListConnectors(r.Context())
		if err != nil {
			handleStorageErrors(w, r, err)
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
