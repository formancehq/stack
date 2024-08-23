package testserver

import (
	"context"
	"encoding/json"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/formancehq/stack/ee/ingester/internal"
	clickhouseconnector "github.com/formancehq/stack/ee/ingester/internal/drivers/clickhouse"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
	"sync"
	"testing"
)

type ClickhouseConnector struct {
	client driver.Conn
	dsn    string
	mu     sync.Mutex
	logger logging.Logger
}

func (h *ClickhouseConnector) initClient() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.client == nil {
		var err error
		h.client, err = clickhouseconnector.OpenDB(h.logger, h.dsn, testing.Verbose())
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *ClickhouseConnector) Clear(ctx context.Context) error {
	if err := h.initClient(); err != nil {
		return err
	}
	return h.client.Exec(ctx, "delete from logs where true")
}

func (h *ClickhouseConnector) ReadMessages(ctx context.Context) ([]ingester.LogWithModule, error) {
	if err := h.initClient(); err != nil {
		return nil, err
	}

	rows, err := h.client.Query(ctx, "select id, shard, module, date, type, data from logs")
	if err != nil {
		return nil, err
	}
	ret := make([]ingester.LogWithModule, 0)
	for rows.Next() {
		var payload string
		newLog := ingester.LogWithModule{}
		if err := rows.Scan(&newLog.ID, &newLog.Shard, &newLog.Module, &newLog.Date, &newLog.Type, &payload); err != nil {
			return nil, errors.Wrap(err, "scanning data from database")
		}
		newLog.Payload = json.RawMessage(payload)

		ret = append(ret, newLog)
	}

	return ret, nil
}

func (h *ClickhouseConnector) Config() map[string]any {
	return map[string]any{
		"dsn": h.dsn,
	}
}

func (h *ClickhouseConnector) Name() string {
	return "clickhouse"
}

var _ Connector = &ClickhouseConnector{}

func NewClickhouseConnector(logger logging.Logger, dsn string) Connector {
	return &ClickhouseConnector{
		dsn:    dsn,
		logger: logger,
	}
}
