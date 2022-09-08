package webhooks

import "time"

type Request struct {
	Date         time.Time `json:"date" bson:"date"`
	ID           string    `json:"id" bson:"id"`
	Config       Config    `json:"config" bson:"config"`
	Payload      string    `json:"payload" bson:"payload"`
	StatusCode   int       `json:"statusCode" bson:"statusCode"`
	Success      bool      `json:"success" bson:"success"`
	RetryAttempt int       `json:"retryAttempt,omitempty" bson:"retryAttempt,omitempty"`
	RetryAfter   time.Time `json:"retryAfter,omitempty" bson:"retryAfter,omitempty"`
}
