package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Log struct {
	bun.BaseModel `bun:"table:logs"`

	ID        string    `json:"id" bun:",pk"`
	Channel   Channel   `json:"channel" bun:"channel"`
	Payload   string    `json:"payload" bun:"payload"`
	CreatedAt time.Time `json:"createdAt" bun:"created_at,nullzero,notnull,default:current_timestamp"`
}

func LogFromEvent(ev Event) (Log, error) {

	payload, err := ev.ToPayload()
	if err != nil {
		return Log{}, err
	}
	return Log{
		ID:        uuid.NewString(),
		Channel:   ev.Channel,
		Payload:   payload,
		CreatedAt: time.Now(),
	}, nil
}
