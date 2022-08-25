package model

import (
	"errors"
	"time"
)

var (
	ErrDateZero     = errors.New("date cannot be zero")
	ErrTypeEmpty    = errors.New("type cannot be empty")
	ErrPayloadEmpty = errors.New("payload cannot be empty")
)

type KafkaEvent struct {
	Date    time.Time      `json:"date"`
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

func (e KafkaEvent) Validate() error {
	if e.Date.IsZero() {
		return ErrDateZero
	}

	if e.Type == "" {
		return ErrTypeEmpty
	}

	if len(e.Payload) == 0 {
		return ErrPayloadEmpty
	}

	return nil
}
