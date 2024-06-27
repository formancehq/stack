package interfaces

import (
	"time"
	"github.com/formancehq/webhooks/internal/commons"
)


type IStoreProvider interface {

	GetHook(index string) (commons.Hook, error)
	SaveHook(hooks commons.Hook) error	
	
	ActivateHook(index string) (commons.Hook, error)
	DeactivateHook(index string) (commons.Hook, error)
	DeleteHook(index string) (commons.Hook, error)

	UpdateHookSecret(index string, secret string) (commons.Hook, error)
	UpdateHookEndpoint(index string, endpoint string) (commons.Hook, error)
	UpdateHookRetry(index string, retry bool) (commons.Hook, error)

	GetHooks(page int, size int, filterEndpoint string) (*[]*commons.Hook, bool,  error)
	LoadHooks() (*[]*commons.Hook, error)


	GetAttempt(index string) (commons.Attempt, error)
	SaveAttempt(attempts commons.Attempt, wrapInLog bool) (error)

	CompleteAttempt(index string) (commons.Attempt, error)
	AbortAttempt(index string, comment string, wrapInLog bool) (commons.Attempt, error)
	ChangeAttemptStatus(index string, status commons.AttemptStatus, comment string, wrapInLog bool) (commons.Attempt, error)
	UpdateAttemptNextTry(index string, nextTry time.Time, statusCode int)(commons.Attempt, error)

	GetWaitingAttempts(page int, size int) (*[]*commons.Attempt, bool, error)
	LoadWaitingAttempts() (*[]*commons.Attempt, error)

	GetAbortedAttempts(page int, size int) (*[]*commons.Attempt, bool, error)

	WriteLog(id, channel, payload string, created_time time.Time) error 
	
	FlushAttempts(index string) (error)
	

	GetFreshLogs(channels []commons.Channel, lastDate time.Time) (*[]*commons.Log, error)
	ListenUpdates(delay int, channels ...commons.Channel) (chan commons.Event, error)
	
	Close() error

}