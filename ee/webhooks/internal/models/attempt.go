package models

import (
	"time"

	"github.com/formancehq/stack/libs/go-libs/sync/shared"
	"github.com/formancehq/webhooks/internal/migrations"
	"github.com/google/uuid"
)

// #############################################################################
// #############################################################################

type AttemptStatus string

const (
	WaitingStatus AttemptStatus = "WAITING"
	SuccessStatus AttemptStatus = "SUCCESS"
	AbortStatus   AttemptStatus = "ABORT"
)

type CommentStatus string

const (
	AbortMissingHook  CommentStatus = "ABORTED FOR 'MISSING OR DELETED HOOK' REASON"
	AbortDisabledHook CommentStatus = "ABORTED FOR 'DISABLED HOOK' REASON"
	AbortNoRetryMode  CommentStatus = "ABORTED FOR 'NO RETRY' REASON"
	AbortMaxRetry     CommentStatus = "ABORTED FOR 'MAX RETRY' REASON"
	AbortUser         CommentStatus = "ABORTED BY USER"
)

type Attempt struct {
	ID                 string        `json:"id" bun:",pk"`
	HookID             string        `json:"hookId" bun:"webhook_id"`
	HookName           string        `json:"hookName" bun:"hook_name"`
	HookEndpoint       string        `json:"hookEndpoint" bun:"hook_endpoint"`
	Event              string        `json:"event" bun:"event"`
	Payload            string        `json:"payload" bun:"payload"`
	LastHttpStatusCode int           `json:"statusCode" bun:"status_code"`
	DateOccured        time.Time     `json:"dateOccured" bun:"date_occured"`
	Status             AttemptStatus `json:"status" bun:"status"`
	DateStatus         time.Time     `json:"dateStatus" bun:"date_status"`
	Comment            CommentStatus `json:"comment" bun:"comment"`
	NextTry            time.Time     `json:"nextRetryAfter,omitempty" bun:"next_retry_after,nullzero"`
	NbTry              int

	CreatedAt    time.Time         `json:"createdAt" bun:"created_at,nullzero,notnull,default:current_timestamp"` //v1
	UpdatedAt    time.Time         `json:"updatedAt" bun:"updated_at,nullzero,notnull,default:current_timestamp"` //v1
	Config       migrations.Config `json:"config" bun:"type:jsonb"`                                               //V1
	RetryAttempt int               `json:"retryAttempt" bun:"retry_attempt"`                                      //V1
}

func NewAttempt(hookID, hookName, hookEndpoint, event, payload string) *Attempt {
	return &Attempt{
		ID:           uuid.NewString(),
		HookID:       hookID,
		HookName:     hookName,
		HookEndpoint: hookEndpoint,
		Event:        event,
		Payload:      payload,
		DateOccured:  time.Now(),
		Status:       WaitingStatus,
		DateStatus:   time.Now(),
		NextTry:      time.Now(),
		NbTry:        0,
	}
}

func (a *Attempt) IsSuccess() bool {
	return a.Status == SuccessStatus
}

func SetAbortMissingHookStatus(a *Attempt) {
	a.Status = AbortStatus
	a.Comment = AbortMissingHook
	a.DateStatus = time.Now()
}

func SetAbortNoRetryModeStatus(a *Attempt) {
	a.Status = AbortStatus
	a.Comment = AbortNoRetryMode
	a.DateStatus = time.Now()
}

func SetAbortMaxRetryStatus(a *Attempt) {
	a.Status = AbortStatus
	a.Comment = AbortMaxRetry
	a.DateStatus = time.Now()
}

func SetAbortUser(a *Attempt) {
	a.Status = AbortStatus
	a.Comment = AbortUser
	a.DateStatus = time.Now()

}

func SetAbortDisableHook(a *Attempt) {
	a.Status = AbortStatus
	a.Comment = AbortDisabledHook
	a.DateStatus = time.Now()
}

func SetSuccesStatus(a *Attempt) {
	a.Status = SuccessStatus
	a.DateStatus = time.Now()
}

func SetNextRetry(a *Attempt) {
	now := time.Now()
	a.NextTry = now.Add(time.Duration(1<<a.NbTry) * time.Second)
}

// #############################################################################
// #############################################################################

type SharedAttempt = shared.Shared[Attempt]

func NewSharedAttempt(hookID, hookName, hookEndpoint, event, payload string) *SharedAttempt {
	s := shared.NewShared(NewAttempt(hookID, hookName, hookEndpoint, event, payload))
	return &s
}

// #############################################################################
// #############################################################################

type SharedAttempts = shared.SharedArr[Attempt]

func NewSharedAttempts() *SharedAttempts {
	s := shared.NewSharedArr[Attempt]()
	return &s
}

// #############################################################################
// #############################################################################

type MapSharedAttempt = shared.SharedMap[Attempt]

func NewMapSharedAttempt() *MapSharedAttempt {
	m := shared.NewSharedMap[Attempt]()
	return &m
}
