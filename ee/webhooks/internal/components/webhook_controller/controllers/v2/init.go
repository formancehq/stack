package controllers

import (
	
	serverInterfaces "github.com/formancehq/webhooks/internal/services/httpserver/interfaces"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
	clientInterface "github.com/formancehq/webhooks/internal/services/httpclient/interfaces"
)


func RegisterV2Controllers(serverHttp serverInterfaces.IHTTPServer, database storeInterface.IStoreProvider, clientHttp clientInterface.IHTTPClient){
	RegisterV2HookControllers(serverHttp, database, clientHttp)
	RegisterV2AttemptControllers(serverHttp, database)
}