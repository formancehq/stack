package logs

import (
	"github.com/formancehq/stack/libs/go-libs/time"
)

type Log struct {
	Date    time.Time `json:"date"`
	Version string    `json:"version"`
	Type    string    `json:"type"`
	Payload any       `json:"payload"`
}
