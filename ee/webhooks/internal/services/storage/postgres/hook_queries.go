package storage

import (
	"context"
	"database/sql"
	"strings"

	"github.com/formancehq/webhooks/internal/models"
)

const (
	selectOneHookQuery             string = "SELECT * FROM configs WHERE id = ?"
	selectHooksQuery                      = "SELECT * FROM configs WHERE status != ?"
	selectHooksWithPaginationQuery        = "SELECT * FROM configs WHERE  status !=  ? AND %condWhere% ORDER BY created_at DESC LIMIT ? OFFSET ?"
	insertHookQuery                       = "INSERT INTO configs (id, name, status, event_types, endpoint, secret, created_at, date_status, retry) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *"
	updateHookStatusQuery                 = "UPDATE configs SET status = ?, date_status = NOW() WHERE id = ? RETURNING *"
	updateHookSecretQuery                 = "UPDATE configs SET secret = ? WHERE id = ? RETURNING *"
	updateHookEndpointQuery               = "UPDATE configs SET endpoint = ? WHERE id = ? RETURNING *"
	updateHookRetryQuery                  = "UPDATE configs SET retry = ? WHERE id = ? RETURNING *"
)

const (
	endpointCriteria = "endpoint = '%val%'"
)

func (store PostgresStore) GetHook(index string) (models.Hook, error) {
	var hook models.Hook

	err := store.db.NewRaw(selectOneHookQuery, index).Scan(context.Background(), &hook)

	if err == sql.ErrNoRows {
		return hook, nil
	}

	return hook, err
}

func (store PostgresStore) SaveHook(hook models.Hook) (models.Hook, error) {
	var savedHook models.Hook

	event, err := models.EventFromType(models.NewHookType, nil, &hook)
	if err != nil {
		return savedHook, err
	}
	log, err := models.LogFromEvent(event)

	wrapQuery := wrapWithLogQuery(insertHookQuery)

	err = store.db.NewRaw(wrapQuery,
		hook.ID,
		hook.Name,
		hook.Status,
		StrArray(hook.Events),
		hook.Endpoint,
		hook.Secret,
		hook.DateCreation,
		hook.DateStatus,
		hook.Retry,
		log.ID,
		log.Channel,
		log.Payload,
		log.CreatedAt,
	).Scan(context.Background(), &savedHook)

	return savedHook, err
}

func (store PostgresStore) ActivateHook(index string) (models.Hook, error) {
	return store.changeHookStatus(index, models.EnableStatus)
}

func (store PostgresStore) DeactivateHook(index string) (models.Hook, error) {
	return store.changeHookStatus(index, models.DisableStatus)
}

func (store PostgresStore) DeleteHook(index string) (models.Hook, error) {
	return store.changeHookStatus(index, models.DeleteStatus)
}

func (store PostgresStore) changeHookStatus(index string, status models.HookStatus) (models.Hook, error) {
	var resHook models.Hook
	var hook models.Hook
	hook.ID = index
	hook.Status = status

	event, err := models.EventFromType(models.ChangeHookStatusType, nil, &hook)
	if err != nil {
		return resHook, err
	}
	log, err := models.LogFromEvent(event)

	wrapQuery := wrapWithLogQuery(updateHookStatusQuery)

	err = store.db.NewRaw(wrapQuery, string(status), index, log.ID, log.Channel, log.Payload, log.CreatedAt).Scan(context.Background(), &resHook)

	if err == sql.ErrNoRows {
		return resHook, nil
	}

	return resHook, err
}

func (store PostgresStore) UpdateHookEndpoint(index string, endpoint string) (models.Hook, error) {
	var resHook models.Hook
	var hook models.Hook
	hook.ID = index
	hook.Endpoint = endpoint

	event, err := models.EventFromType(models.ChangeHookEndpointType, nil, &hook)
	if err != nil {
		return resHook, err
	}
	log, err := models.LogFromEvent(event)

	wrapQuery := wrapWithLogQuery(updateHookEndpointQuery)

	_, err = store.db.NewRaw(wrapQuery, string(endpoint), index, log.ID, log.Channel, log.Payload, log.CreatedAt).Exec(context.Background(), &resHook)

	if err == sql.ErrNoRows {
		return resHook, nil
	}

	return resHook, err
}

func (store PostgresStore) UpdateHookSecret(index string, secret string) (models.Hook, error) {
	var resHook models.Hook
	var hook models.Hook
	hook.ID = index
	hook.Secret = secret

	event, err := models.EventFromType(models.ChangeHookSecretType, nil, &hook)
	if err != nil {
		return resHook, err
	}
	log, err := models.LogFromEvent(event)

	wrapQuery := wrapWithLogQuery(updateHookSecretQuery)

	_, err = store.db.NewRaw(wrapQuery, secret, index, log.ID, log.Channel, log.Payload, log.CreatedAt).Exec(context.Background(), &resHook)

	if err == sql.ErrNoRows {
		return resHook, nil
	}

	return resHook, err

}

func (store PostgresStore) UpdateHookRetry(index string, retry bool) (models.Hook, error) {
	var resHook models.Hook
	var hook models.Hook
	hook.ID = index
	hook.Retry = retry

	event, err := models.EventFromType(models.ChangeHookRetryType, nil, &hook)
	if err != nil {
		return resHook, err
	}
	log, err := models.LogFromEvent(event)

	wrapQuery := wrapWithLogQuery(updateHookRetryQuery)

	_, err = store.db.NewRaw(wrapQuery, retry, index, log.ID, log.Channel, log.Payload, log.CreatedAt).Exec(context.Background(), &resHook)

	if err == sql.ErrNoRows {
		return resHook, nil
	}

	return resHook, err

}

func (store PostgresStore) GetHooks(page, size int, filterEndpoint string) (*[]*models.Hook, bool, error) {
	res := make([]*models.Hook, 0)
	hasMore := false

	condWhere := "1=1"

	rawQuery := selectHooksWithPaginationQuery
	if filterEndpoint != "" {
		condWhere = strings.Replace(endpointCriteria, "%val%", filterEndpoint, 1)
	}

	rawQuery = strings.Replace(rawQuery, "%condWhere%", condWhere, 1)

	err := store.db.NewRaw(rawQuery, models.DeleteStatus, size+1, size*page).Scan(context.Background(), &res)

	if err != nil {
		return &res, hasMore, err
	}

	hasMore = len(res) == (size + 1)

	if hasMore {
		res = res[0:size]
	}

	return &res, hasMore, err

}

func (store PostgresStore) LoadHooks() (*[]*models.Hook, error) {
	res := make([]*models.Hook, 0)

	_, err := store.db.NewRaw(selectHooksQuery, models.DeleteStatus).Exec(context.Background(), &res)

	return &res, err
}
