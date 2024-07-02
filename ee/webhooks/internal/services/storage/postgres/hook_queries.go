package storage

import (
	"context"
	"database/sql"
	"strings"

	"github.com/formancehq/webhooks/internal/commons"
)

const (
	selectOneHookQuery             string = "SELECT * FROM configs WHERE id = ?"
	selectHooksQuery                      = "SELECT * FROM configs WHERE status != ?"
	selectHooksWithPaginationQuery        = "SELECT * FROM configs WHERE  status !=  ? AND %condWhere% ORDER By name LIMIT ? OFFSET ?"
	insertHookQuery                       = "INSERT INTO configs (id, name, status, event_types, endpoint, secret, created_at, date_status, retry) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *"
	updateHookStatusQuery                 = "UPDATE configs SET status = ?, date_status = NOW() WHERE id = ? RETURNING *"
	updateHookSecretQuery                 = "UPDATE configs SET secret = ? WHERE id = ? RETURNING *"
	updateHookEndpointQuery               = "UPDATE configs SET endpoint = ? WHERE id = ? RETURNING *"
	updateHookRetryQuery                  = "UPDATE configs SET retry = ? WHERE id = ? RETURNING *"
)

const (
	endpointCriteria = "endpoint = '%val%'"
)

func (store PostgresStore) GetHook(index string) (commons.Hook, error) {
	var hook commons.Hook

	err := store.db.NewRaw(selectOneHookQuery, index).Scan(context.Background(), &hook)

	if err == sql.ErrNoRows {
		return hook, nil
	}

	return hook, err
}

func (store PostgresStore) SaveHook(hook commons.Hook) error {

	event, err := commons.EventFromType(commons.NewHookType, nil, &hook)
	if err != nil {
		return err
	}
	log, err := commons.LogFromEvent(event)

	wrapQuery := wrapWithLogQuery(insertHookQuery)

	_, err = store.db.NewRaw(wrapQuery,
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
	).
		Exec(context.Background())

	return err
}

func (store PostgresStore) ActivateHook(index string) (commons.Hook, error) {
	return store.changeHookStatus(index, commons.EnableStatus)
}

func (store PostgresStore) DeactivateHook(index string) (commons.Hook, error) {
	return store.changeHookStatus(index, commons.DisableStatus)
}

func (store PostgresStore) DeleteHook(index string) (commons.Hook, error) {
	return store.changeHookStatus(index, commons.DeleteStatus)
}

func (store PostgresStore) changeHookStatus(index string, status commons.HookStatus) (commons.Hook, error) {
	var hook commons.Hook
	hook.ID = index
	hook.Status = status

	event, err := commons.EventFromType(commons.ChangeHookStatusType, nil, &hook)
	if err != nil {
		return hook, err
	}
	log, err := commons.LogFromEvent(event)

	wrapQuery := wrapWithLogQuery(updateHookStatusQuery)

	_, err = store.db.NewRaw(wrapQuery, string(status), index, log.ID, log.Channel, log.Payload, log.CreatedAt).Exec(context.Background(), &hook)
	_, err = store.db.NewRaw(wrapQuery, string(status), index, log.ID, log.Channel, log.Payload, log.CreatedAt).Exec(context.Background(), &hook)

	if err == sql.ErrNoRows {
		hook.ID = ""
		return hook, nil
	}

	return hook, err
}

func (store PostgresStore) UpdateHookEndpoint(index string, endpoint string) (commons.Hook, error) {
	var hook commons.Hook
	hook.ID = index
	hook.Endpoint = endpoint

	event, err := commons.EventFromType(commons.ChangeHookEndpointType, nil, &hook)
	if err != nil {
		return hook, err
	}
	log, err := commons.LogFromEvent(event)

	wrapQuery := wrapWithLogQuery(updateHookEndpointQuery)

	_, err = store.db.NewRaw(wrapQuery, string(endpoint), index, log.ID, log.Channel, log.Payload, log.CreatedAt).Exec(context.Background(), &hook)

	if err == sql.ErrNoRows {
		hook.ID = ""
		return hook, nil
	}

	return hook, err
}

func (store PostgresStore) UpdateHookSecret(index string, secret string) (commons.Hook, error) {
	var hook commons.Hook
	hook.ID = index
	hook.Secret = secret

	event, err := commons.EventFromType(commons.ChangeHookSecretType, nil, &hook)
	if err != nil {
		return hook, err
	}
	log, err := commons.LogFromEvent(event)

	wrapQuery := wrapWithLogQuery(updateHookSecretQuery)

	_, err = store.db.NewRaw(wrapQuery, secret, index, log.ID, log.Channel, log.Payload, log.CreatedAt).Exec(context.Background(), &hook)

	if err == sql.ErrNoRows {
		hook.ID = ""
		return hook, nil
	}

	return hook, err

}

func (store PostgresStore) UpdateHookRetry(index string, retry bool) (commons.Hook, error) {
	var hook commons.Hook
	hook.ID = index
	hook.Retry = retry

	event, err := commons.EventFromType(commons.ChangeHookRetryType, nil, &hook)
	if err != nil {
		return hook, err
	}
	log, err := commons.LogFromEvent(event)

	wrapQuery := wrapWithLogQuery(updateHookRetryQuery)

	_, err = store.db.NewRaw(wrapQuery, retry, index, log.ID, log.Channel, log.Payload, log.CreatedAt).Exec(context.Background(), &hook)

	if err == sql.ErrNoRows {
		hook.ID = ""
		return hook, nil
	}

	return hook, err

}

func (store PostgresStore) GetHooks(page, size int, filterEndpoint string) (*[]*commons.Hook, bool, error) {
	res := make([]*commons.Hook, 0)
	hasMore := false

	condWhere := "1=1"

	rawQuery := selectHooksWithPaginationQuery
	if filterEndpoint != "" {
		condWhere = strings.Replace(endpointCriteria, "%val%", filterEndpoint, 1)
	}

	rawQuery = strings.Replace(rawQuery, "%condWhere%", condWhere, 1)

	_, err := store.db.NewRaw(rawQuery, commons.DeleteStatus, size+1, size*page).Exec(context.Background(), &res)

	if err != nil {
		return &res, hasMore, err
	}

	hasMore = len(res) == (size + 1)

	if hasMore {
		res = res[0:size]
	}

	return &res, hasMore, err

}

func (store PostgresStore) LoadHooks() (*[]*commons.Hook, error) {
	res := make([]*commons.Hook, 0)

	_, err := store.db.NewRaw(selectHooksQuery, commons.DeleteStatus).Exec(context.Background(), &res)

	return &res, err
}
