package messages

import "time"

const (
	TopicAudit = "payments"

	EventVersion = "v1"
	EventApp     = "gateway"

	EventTypeAudit = "AUDIT"
)

type EventMessage struct {
	Date    time.Time `json:"date"`
	App     string    `json:"app"`
	Version string    `json:"version"`
	Type    string    `json:"type"`
	Payload any       `json:"payload"`
}
