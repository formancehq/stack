package utils

import (
	"time"

	"github.com/formancehq/webhooks/internal/models"
)

type V1HookUser struct {
	Endpoint   string   `json:"endpoint"`
	Secret     string   `json:"secret"`
	EventTypes []string `json:"eventTypes" `
}

type V1Hook struct {
	V1HookUser

	ID        string    `json:"id"`
	Active    bool      `json:"active"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt" `
	UpdatedAt time.Time `json:"updatedAt" `
}

type V1Attempt struct {
	ID             string    `json:"id"`
	WebhookID      string    `json:"webhookID"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Config         V1Hook    `json:"config"`
	Payload        string    `json:"payload"`
	StatusCode     int       `json:"statusCode"`
	RetryAttempt   int       `json:"retryAttempt"`
	Status         string    `json:"status"`
	NextRetryAfter time.Time `json:"nextRetryAfter,omitempty" bun:"next_retry_after,nullzero"`
}

func ToV1Hook(hook models.Hook) V1Hook {

	c := V1HookUser{
		Endpoint:   hook.Endpoint,
		Secret:     hook.Secret,
		EventTypes: hook.Events,
	}

	return V1Hook{
		V1HookUser: c,
		ID:         hook.ID,
		Active:     hook.IsActive(),
		Name:       hook.Name,
		CreatedAt:  hook.DateCreation,
		UpdatedAt:  hook.DateStatus,
	}

}

func ToV1Hooks(hooks *[]*models.Hook) []V1Hook {
	v1Hooks := make([]V1Hook, 0)
	for _, hook := range *hooks {
		v1Hooks = append(v1Hooks, ToV1Hook(*hook))
	}
	return v1Hooks
}

func ToV1Attempt(hook models.Hook, attempt models.Attempt) V1Attempt {

	var status string
	if attempt.Status == models.SuccessStatus {
		status = "success"
	} else {
		status = "failed"
	}
	return V1Attempt{
		ID:             attempt.ID,
		WebhookID:      hook.ID,
		Config:         ToV1Hook(hook),
		Payload:        attempt.Payload,
		StatusCode:     attempt.LastHttpStatusCode,
		Status:         status,
		RetryAttempt:   attempt.RetryAttempt,
		UpdatedAt:      attempt.DateStatus,
		NextRetryAfter: attempt.NextTry,
		CreatedAt:      attempt.CreatedAt,
	}
}
