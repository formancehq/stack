package interfaces

import (
	"time"

	"github.com/formancehq/webhooks/internal/models"
)

type IStoreProvider interface {
	GetHook(index string) (models.Hook, error)
	SaveHook(hooks models.Hook) (models.Hook, error)

	ActivateHook(index string) (models.Hook, error)
	DeactivateHook(index string) (models.Hook, error)
	DeleteHook(index string) (models.Hook, error)

	UpdateHookSecret(index string, secret string) (models.Hook, error)
	UpdateHookEndpoint(index string, endpoint string) (models.Hook, error)
	UpdateHookRetry(index string, retry bool) (models.Hook, error)

	GetHooks(page int, size int, filterEndpoint string) (*[]*models.Hook, bool, error)
	LoadHooks() (*[]*models.Hook, error)

	GetAttempt(index string) (models.Attempt, error)
	SaveAttempt(attempts models.Attempt, wrapInLog bool) error

	CompleteAttempt(index string) (models.Attempt, error)
	AbortAttempt(index string, comment string, wrapInLog bool) (models.Attempt, error)
	ChangeAttemptStatus(index string, status models.AttemptStatus, comment string, wrapInLog bool) (models.Attempt, error)
	UpdateAttemptEndpoint(index string, endpoint string) (models.Attempt, error)
	UpdateAttemptNextTry(index string, nextTry time.Time, statusCode int) (models.Attempt, error)

	GetWaitingAttempts(page int, size int) (*[]*models.Attempt, bool, error)
	LoadWaitingAttempts() (*[]*models.Attempt, error)

	GetAbortedAttempts(page int, size int) (*[]*models.Attempt, bool, error)

	WriteLog(id, channel, payload string, created_time time.Time) error

	FlushAttempts(index string) error

	GetFreshLogs(channels []models.Channel, lastDate time.Time) (*[]*models.Log, error)
	ListenUpdates(delay int, channels ...models.Channel) (chan models.Event, error)

	Close() error
}
