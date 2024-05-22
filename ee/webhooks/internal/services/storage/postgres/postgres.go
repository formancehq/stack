package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/webhooks/internal/models"
	"github.com/formancehq/webhooks/internal/services/storage/interfaces"
	"github.com/uptrace/bun"
)

type PostgresStore struct {
	db *bun.DB
}

type Query string

func (q Query) Fill(vals ...any) Query {
	return Query(fmt.Sprintf(string(q), vals...))
}

func (q Query) ValuesNb(nb int) string {
	parts := make([]string, nb)
	for i := 0; i < nb; i++ {
		parts[i] = "?"
	}

	return strings.Join(parts, ",")
}

var ChannelHook string = "update_hook"
var ChannelAttempt string = "update_attempt"

func (store PostgresStore) ListenUpdates(delay int, channels ...models.Channel) (chan models.Event, error) {

	if store.db == nil {
		return nil, errors.New("No Database in PostgresStore")
	}
	err := store.db.Ping()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to ping the database: %v", err))
	}
	dbEventChan := make(chan models.Event)
	var lastUpdate time.Time
	lastUpdate = time.Now()

	go func() {
		for {
			select {
			case <-time.After(time.Duration(delay) * time.Second):
				newLastUpdate := time.Now()
				query := store.db.NewSelect().Table("logs").ColumnExpr("*").Order("created_at ASC").
					Where("channel IN (?)", bun.In(channels)).
					Where("created_at > ?", lastUpdate)

				logs := make([]models.Log, 0)
				err := query.Scan(context.Background(), &logs)
				if err != nil {

					if err != nil {
						message := fmt.Sprintf("PostgresStore:ListenUpdate() - goroutine :  Error while attempting to reach the database and read logs: %s", err)
						logging.Error(message)
						panic(message)
					}
				}
				for _, log := range logs {
					event, err := models.Event{}.FromPayload(log.Payload)
					if err != nil {
						message := fmt.Sprintf("PostgresStore:ListenUpdate() - goroutine :  Error while attempt to read payload from logs: %s", err)
						logging.Error(message)
						panic(message)
					}
					dbEventChan <- event

				}
				lastUpdate = newLastUpdate
			}
		}
	}()

	return dbEventChan, nil
}

func (store PostgresStore) Close() error {
	return store.db.Close()
}

func NewPostgresStoreProvider(db *bun.DB) PostgresStore {
	postgresStore := PostgresStore{
		db: db,
	}

	return postgresStore
}

var _ interfaces.IStoreProvider = (*PostgresStore)(nil)
