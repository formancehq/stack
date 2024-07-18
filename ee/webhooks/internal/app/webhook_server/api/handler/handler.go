package handler

import (
	businessService "github.com/formancehq/webhooks/internal/app/webhook_server/api/service"
	clientInterface "github.com/formancehq/webhooks/internal/services/httpclient/interfaces"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
)

const (
	hookPageSize    int = 80
	attemptPageSize int = 80
)

type PayloadBody struct {
	Payload string `json:"payload"`
}

func SetDatabase(db storeInterface.IStoreProvider) {

	businessService.SetDatabase(db)
}

func SetClientHTTP(c clientInterface.IHTTPClient) {
	businessService.SetClientHTTP(c)
}
