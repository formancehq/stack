package commons

import(
	"github.com/formancehq/webhooks/internal/commons"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
)

func SendEvent(database storeInterface.IStoreProvider, t commons.EventType, attempt *commons.Attempt, hook *commons.Hook) error{
	event, err := commons.EventFromType(t, attempt, hook)
	if err != nil {return err}
	
	return database.NotifyUpdate(event) 
}