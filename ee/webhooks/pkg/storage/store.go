package storage

import (
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/pkg/errors"
)

var (
	ErrConfigNotFound      = errors.New("config not found")
	ErrConfigNotModified   = errors.New("config not modified")
	ErrWebhookIDNotFound   = errors.New("webhook ID not found")
	ErrAttemptsNotModified = errors.New("attempt not modified")
)

type Store interface {
	webhooks.ConfigStore
	webhooks.AttemptStore
}
