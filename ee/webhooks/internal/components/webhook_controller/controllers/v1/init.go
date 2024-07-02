package controllers

import (
	clientInterface "github.com/formancehq/webhooks/internal/services/httpclient/interfaces"
	serverInterfaces "github.com/formancehq/webhooks/internal/services/httpserver/interfaces"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
)


func RegisterV1Controllers(serverHttp serverInterfaces.IHTTPServer, database storeInterface.IStoreProvider, client clientInterface.IHTTPClient) {

	RegisterV1HookControllers(serverHttp, database, client)

}
