package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/formancehq/webhooks/internal/models"
)

const (
	selectOneAttemptQuery             string = "SELECT * FROM attempts WHERE id = ? "
	selectAttemptsQuery                      = "SELECT * FROM attempts WHERE (status = ? ) OR (status = ?)"
	selectAttemptsWithPaginationQuery string = "SELECT * FROM attempts WHERE (status = ? ) OR (status = ?) ORDER BY date_status DESC LIMIT ? OFFSET ?"
	insertAttemptQuery                       = "INSERT INTO attempts (id, webhook_id, hook_name, hook_endpoint, event, payload, status_code, date_occured, status, date_status, comment, next_retry_after) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *"
	updateAttemptStatus                      = "UPDATE attempts SET status = ?, date_status = NOW(), comment = ? WHERE id = ? RETURNING *"
	updateAttemptNextTry                     = "UPDATE attempts SET next_retry_after = ?, status_code = ? WHERE id = ? RETURNING *"
	updateAttemptEndpoint                    = "UPDATE attempts SET hook_endpoint = ? WHERE id = ? RETURNING *"
)

func (store PostgresStore) GetAttempt(index string) (models.Attempt, error) {
	var attempt models.Attempt

	err := store.db.NewRaw(selectOneAttemptQuery, index).Scan(context.Background(), &attempt)

	if err == sql.ErrNoRows {
		return attempt, nil
	}

	return attempt, err

}

func (store PostgresStore) SaveAttempt(attempt models.Attempt, wrapInLog bool) error {

	var err error

	if wrapInLog {
		event, err := models.EventFromType(models.NewWaitingAttemptType, &attempt, nil)
		if err != nil {
			return err
		}
		log, err := models.LogFromEvent(event)

		wrapQuery := wrapWithLogQuery(insertAttemptQuery)

		_, err = store.db.NewRaw(wrapQuery,
			attempt.ID,
			attempt.HookID,
			attempt.HookName,
			attempt.HookEndpoint,
			attempt.Event,
			attempt.Payload,
			attempt.LastHttpStatusCode,
			attempt.DateOccured,
			string(attempt.Status),
			attempt.DateStatus,
			string(attempt.Comment),
			attempt.NextTry,
			log.ID,
			log.Channel,
			log.Payload,
			log.CreatedAt).
			Exec(context.Background())

	} else {
		_, err = store.db.NewRaw(insertAttemptQuery,
			attempt.ID,
			attempt.HookID,
			attempt.HookName,
			attempt.HookEndpoint,
			attempt.Event,
			attempt.Payload,
			attempt.LastHttpStatusCode,
			attempt.DateOccured,
			string(attempt.Status),
			attempt.DateStatus,
			string(attempt.Comment),
			attempt.NextTry).
			Exec(context.Background())
	}

	return err
}
func (store PostgresStore) CompleteAttempt(index string) (models.Attempt, error) {
	return store.ChangeAttemptStatus(index, models.SuccessStatus, "", false)
}

func (store PostgresStore) AbortAttempt(index string, comment string, wrapInLog bool) (models.Attempt, error) {
	return store.ChangeAttemptStatus(index, models.AbortStatus, comment, wrapInLog)
}

func (store PostgresStore) ChangeAttemptStatus(index string, status models.AttemptStatus, comment string, wrapInLog bool) (models.Attempt, error) {
	var attempt models.Attempt
	attempt.ID = index

	var err error

	if wrapInLog {

		event, err := models.EventFromType(models.AbortWaitingAttemptType, &attempt, nil)
		if err != nil {
			return attempt, err
		}
		log, err := models.LogFromEvent(event)

		wrapQuery := wrapWithLogQuery(updateAttemptStatus)
		query := store.db.NewRaw(wrapQuery,
			string(status),
			comment,
			index,
			log.ID,
			log.Channel,
			log.Payload,
			log.CreatedAt,
		)

		err = query.Scan(context.Background(), &attempt)

	} else {
		query := store.db.NewRaw(updateAttemptStatus,
			string(status),
			comment,
			index,
		)

		err = query.Scan(context.Background(), &attempt)
	}

	if err == sql.ErrNoRows {
		attempt.ID = ""
		return attempt, err
	}

	return attempt, err
}

func (store PostgresStore) UpdateAttemptEndpoint(index string, endpoint string) (models.Attempt, error) {
	var attempt models.Attempt

	_, err := store.db.NewRaw(updateAttemptEndpoint, endpoint, index).Exec(context.Background(), &attempt)

	if err == sql.ErrNoRows {
		return attempt, nil
	}

	return attempt, nil
}

func (store PostgresStore) UpdateAttemptNextTry(index string, nextTry time.Time, statusCode int) (models.Attempt, error) {
	var attempt models.Attempt

	_, err := store.db.NewRaw(updateAttemptNextTry, nextTry, statusCode, index).Exec(context.Background(), &attempt)

	if err == sql.ErrNoRows {
		return attempt, nil
	}

	return attempt, err
}

func (store PostgresStore) GetWaitingAttempts(page, size int) (*[]*models.Attempt, bool, error) {
	res := make([]*models.Attempt, 0)
	hasMore := false

	err := store.db.NewRaw(selectAttemptsWithPaginationQuery, models.WaitingStatus, "to_retry", size+1, size*page).
		Scan(context.Background(), &res)

	hasMore = len(res) == (size + 1)

	if hasMore {
		res = res[0:size]
	}

	return &res, hasMore, err
}

func (store PostgresStore) GetAbortedAttempts(page, size int) (*[]*models.Attempt, bool, error) {
	res := make([]*models.Attempt, 0)
	hasMore := false

	err := store.db.NewRaw(selectAttemptsWithPaginationQuery, models.AbortStatus, "failed", size+1, size*page).
		Scan(context.Background(), &res)

	hasMore = len(res) == (size + 1)

	if hasMore {
		res = res[0:size]
	}

	return &res, hasMore, err
}

func (store PostgresStore) LoadWaitingAttempts() (*[]*models.Attempt, error) {
	res := make([]*models.Attempt, 0)

	err := store.db.NewRaw(selectAttemptsQuery, models.WaitingStatus, "to_retry").Scan(context.Background(), &res)

	return &res, err
}

func (store PostgresStore) FlushAttempts(index string) error {

	attempt := models.Attempt{}
	var log models.Log
	var err error

	if index != "" {

		attempt.ID = index
		event, err := models.EventFromType(models.FlushWaitingAttemptType, &attempt, nil)
		if err != nil {
			return err
		}

		log, err = models.LogFromEvent(event)
		if err != nil {
			return err
		}

	} else {

		event, err := models.EventFromType(models.FlushWaitingAttemptsType, &attempt, nil)
		if err != nil {
			return err
		}

		log, err = models.LogFromEvent(event)
		if err != nil {
			return err
		}

	}

	err = store.db.NewRaw(insertLogQuery, log.ID, log.Channel, log.Payload, log.CreatedAt).Scan(context.Background())

	return err
}
