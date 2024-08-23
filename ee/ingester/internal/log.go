package ingester

import (
	"encoding/json"
	"time"
)

type Log struct {
	Shard   string          `json:"shard"`
	ID      string          `json:"id"`
	Date    time.Time       `json:"date"`
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type LogWithModule struct {
	Log
	Module string
}

func NewLogWithModule(module string, log Log) LogWithModule {
	return LogWithModule{
		Log:    log,
		Module: module,
	}
}
