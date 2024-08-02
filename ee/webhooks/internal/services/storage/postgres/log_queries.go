package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/webhooks/internal/models"
	"github.com/uptrace/bun"
)

const (
	insertLogQuery  string = "INSERT INTO logs (id, channel, payload, created_at) VALUES (?, ?, ?, ?);"
	selectFreshLogs string = "SELECT * FROM logs WHERE channel IN (?) AND created_at > ? ORDER BY created_at ASC;"
)

func wrapWithLogQuery(query string) string {

	rawQuery := `
		BEGIN ;

		WITH wrap_query AS (
			%s
		),

		inserted_log AS (
			INSERT INTO logs (id, channel, payload, created_at)
			VALUES (?, ?, ?, ?)
		)

		SELECT * FROM wrap_query;
		
		COMMIT; 
	`

	return fmt.Sprintf(rawQuery, query)

}

func (store PostgresStore) WriteLog(id, channel, payload string, created_at time.Time) error {

	_, err := store.db.NewRaw(insertLogQuery, id, channel, payload, created_at).Exec(context.Background())

	return err

}

func (store PostgresStore) GetFreshLogs(channels []models.Channel, lastDate time.Time) (*[]*models.Log, error) {
	res := make([]*models.Log, 0)
	_, err := store.db.NewRaw(selectFreshLogs, bun.In(channels), lastDate).Exec(context.Background(), &res)

	return &res, err
}
