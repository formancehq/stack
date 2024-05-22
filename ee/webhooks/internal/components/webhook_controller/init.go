package webhookcontroller

import (

	V1Controllers "github.com/formancehq/webhooks/internal/components/webhook_controller/controllers/v1"
	V2Controllers "github.com/formancehq/webhooks/internal/components/webhook_controller/controllers/v2"

	serverInterfaces "github.com/formancehq/webhooks/internal/services/httpserver/interfaces"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
	clientInterface "github.com/formancehq/webhooks/internal/services/httpclient/interfaces"
)


func Init(serverHttp serverInterfaces.IHTTPServer, database storeInterface.IStoreProvider, client clientInterface.IHTTPClient){
	
	V1Controllers.RegisterV1Controllers(serverHttp, database, client)
	V2Controllers.RegisterV2Controllers(serverHttp, database, client)


}
