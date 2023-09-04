package ledgerstore

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/storage/paginate"
	"github.com/lib/pq"
	"github.com/uptrace/bun"
)

var (
	batchSize             uint64 = 10000
	oldSchemaRenameSuffix        = "_save_v2_0_0"
)

type LogV1 struct {
	bun.BaseModel `bun:"log,alias:log"`

	ID   uint64          `bun:"id,unique,type:bigint"`
	Type string          `bun:"type,type:varchar"`
	Hash string          `bun:"hash,type:varchar"`
	Date ledger.Time     `bun:"date,type:timestamptz"`
	Data json.RawMessage `bun:"data,type:jsonb"`
}

func readLogsRange(
	ctx context.Context,
	schema string,
	sqlTx bun.Tx,
	idMin, idMax uint64,
) ([]LogV1, error) {
	rawLogs := make([]LogV1, 0)
	if err := sqlTx.
		NewSelect().
		Table(fmt.Sprintf(`"%s".log`, schema)).
		Where("id >= ?", idMin).
		Where("id < ?", idMax).
		Model(&rawLogs).
		Scan(ctx); err != nil {
		return nil, err
	}

	return rawLogs, nil
}

func convertMetadata(data []byte) any {
	ret := make(map[string]any)
	if err := json.Unmarshal(data, &ret); err != nil {
		panic(err)
	}
	oldMetadata := ret["metadata"].(map[string]any)
	newMetadata := make(map[string]string)
	for k, v := range oldMetadata {
		newMetadata[k] = fmt.Sprint(v)
	}
	ret["metadata"] = newMetadata

	return ret
}

func (l *LogV1) ToLogsV2() (Logs, error) {
	logType := ledger.LogTypeFromString(l.Type)

	var data any
	switch logType {
	case ledger.NewTransactionLogType:
		data = map[string]any{
			"transaction":     convertMetadata(l.Data),
			"accountMetadata": map[string]any{},
		}
	case ledger.SetMetadataLogType:
		data = convertMetadata(l.Data)
	case ledger.RevertedTransactionLogType:
		data = l.Data
	default:
		panic("unknown type " + logType.String())
	}

	asJson, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return Logs{
		ID:   (*paginate.BigInt)(big.NewInt(int64(l.ID))),
		Type: logType.String(),
		Hash: []byte(l.Hash),
		Date: l.Date,
		Data: asJson,
	}, nil
}

func batchLogs(
	schema string,
	sqlTx bun.Tx,
	logs []Logs,
) error {
	// Beware: COPY query is not supported by bun if the pgx driver is used.
	stmt, err := sqlTx.Prepare(pq.CopyInSchema(
		schema,
		"log",
		"id", "type", "hash", "date", "data",
	))
	if err != nil {
		return err
	}

	for _, l := range logs {
		_, err = stmt.Exec(l.ID, l.Type, l.Hash, l.Date, RawMessage(l.Data))
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func migrateLogs(
	ctx context.Context,
	schemaV1Name string,
	schemaV2Name string,
	sqlTx bun.Tx,
) error {

	var idMin uint64
	var idMax = idMin + batchSize
	for {
		logs, err := readLogsRange(ctx, schemaV1Name, sqlTx, idMin, idMax)
		if err != nil {
			return err
		}

		if len(logs) == 0 {
			break
		}

		logsV2 := make([]Logs, 0, len(logs))
		for _, l := range logs {
			logV2, err := l.ToLogsV2()
			if err != nil {
				return err
			}

			logsV2 = append(logsV2, logV2)
		}

		err = batchLogs(schemaV2Name, sqlTx, logsV2)
		if err != nil {
			return err
		}

		idMin = idMax
		idMax = idMin + batchSize
	}

	return nil
}
