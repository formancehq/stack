package service

import (
	clientInterface "github.com/formancehq/webhooks/internal/services/httpclient/interfaces"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
)

var (
	database storeInterface.IStoreProvider
	client   clientInterface.IHTTPClient
)

func SetDatabase(db storeInterface.IStoreProvider) {
	database = db
}

func getDatabase() storeInterface.IStoreProvider {
	return database
}

func SetClientHTTP(c clientInterface.IHTTPClient) {
	client = c
}

func getClient() clientInterface.IHTTPClient {
	return client
}
