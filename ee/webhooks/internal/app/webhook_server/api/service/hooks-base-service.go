package service

import (
	"context"
	"net/http"

	"github.com/formancehq/webhooks/internal/models"
)

func BaseCreateHook(name string, events []string, endpoint string, secret string, retry bool) (models.Hook, error) {
	hook := models.NewHook(name, events, endpoint, secret, retry)
	savedHook, err := getDatabase().SaveHook(*hook)
	return savedHook, err
}

func BaseGetHooks(filterEndpoint string, page int, pageSize int) (*[]*models.Hook, bool, error) {

	hooks, hasMore, err := getDatabase().GetHooks(page, pageSize, filterEndpoint)

	return hooks, hasMore, err
}

func BaseGetHook(id string) (models.Hook, error) {
	return getDatabase().GetHook(id)
}

func BaseDeleteHook(id string) (models.Hook, error) {
	return getDatabase().DeleteHook(id)
}

func BaseActivateHook(id string) (models.Hook, error) {
	return getDatabase().ActivateHook(id)
}

func BaseDeactivateHook(id string) (models.Hook, error) {
	return getDatabase().DeactivateHook(id)
}

func BaseUpdateSecret(id string, secret string) (models.Hook, error) {
	return getDatabase().UpdateHookSecret(id, secret)
}

func BaseUpdateEndpoint(id string, endpoint string) (models.Hook, error) {
	return getDatabase().UpdateHookEndpoint(id, endpoint)
}

func BaseUpdateRetry(id string, retry bool) (models.Hook, error) {
	return getDatabase().UpdateHookRetry(id, retry)
}

func BaseTestHook(id string, payload string) (*models.Hook, *models.Attempt, error) {

	hook, err := database.GetHook(id)
	if err != nil {
		return nil, nil, err
	}
	if hook.ID == "" {
		return &models.Hook{}, &models.Attempt{}, nil
	}

	attempt := models.NewAttempt(hook.ID, hook.Name, hook.Endpoint, hook.Events[0], payload)

	statusCode, err := getClient().Call(context.Background(), &hook, attempt, true)

	if err == nil {
		attempt.LastHttpStatusCode = statusCode
	}

	if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
		models.SetSuccesStatus(attempt)
	} else {
		models.SetAbortMaxRetryStatus(attempt)
	}

	return &hook, attempt, err
}
