package kafka

import (
	"context"

	kafkago "github.com/segmentio/kafka-go"
)

type Reader interface {
	FetchMessage(ctx context.Context) (kafkago.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafkago.Message) error
	Close() error
}
