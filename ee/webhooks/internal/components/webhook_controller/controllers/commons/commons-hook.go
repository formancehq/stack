package commons

import (
	"context"
	"github.com/formancehq/webhooks/internal/commons"
	clientInterface "github.com/formancehq/webhooks/internal/services/httpclient/interfaces"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
	"net/http"
)

func CreateHook(database storeInterface.IStoreProvider, name string, events []string, endpoint string, secret string, retry bool) (commons.Hook, error) {
	hook := commons.NewHook(name, events, endpoint, secret, retry)
	err := database.SaveHook(*hook)
	return *hook, err
}

func GetHooks(database storeInterface.IStoreProvider, filterEndpoint string, page int, pageSize int) (*[]*commons.Hook, bool, error) {

	hooks, hasMore, err := database.GetHooks(page, pageSize, filterEndpoint)

	return hooks, hasMore, err
}

func GetHook(database storeInterface.IStoreProvider, id string) (commons.Hook, error) {
	return database.GetHook(id)
}

func DeleteHook(database storeInterface.IStoreProvider, id string) (commons.Hook, error) {
	return database.DeleteHook(id)
}

func ActivateHook(database storeInterface.IStoreProvider, id string) (commons.Hook, error) {
	return database.ActivateHook(id)
}

func DeactivateHook(database storeInterface.IStoreProvider, id string) (commons.Hook, error) {
	return database.DeactivateHook(id)
}

func UpdateSecret(database storeInterface.IStoreProvider, id string, secret string) (commons.Hook, error) {
	return database.UpdateHookSecret(id, secret)
}

func UpdateEndpoint(database storeInterface.IStoreProvider, id string, endpoint string) (commons.Hook, error) {
	return database.UpdateHookEndpoint(id, endpoint)
}

func UpdateRetry(database storeInterface.IStoreProvider, id string, retry bool) (commons.Hook, error) {
	return database.UpdateHookRetry(id, retry)
}

func TestHook(database storeInterface.IStoreProvider, clienthttp clientInterface.IHTTPClient, id string, payload string) (*commons.Hook, *commons.Attempt, error) {

	hook, err := database.GetHook(id)
	if err != nil {
		return nil, nil, err
	}
	if hook.ID == "" {
		return &commons.Hook{}, &commons.Attempt{}, nil
	}

	attempt := commons.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], payload)

	statusCode, err := clienthttp.Call(context.Background(), &hook, attempt, true)

	if err == nil {
		attempt.LastHttpStatusCode = statusCode
	}

	if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
		commons.SetSuccesStatus(attempt)
	} else {
		commons.SetAbortMaxRetryStatus(attempt)
	}

	return &hook, attempt, err
}
