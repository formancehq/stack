package models

import (
	"time"
)

type Instance struct {
	ID           string
	ScheduleID   string
	ConnectorID  ConnectorID
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Terminated   bool
	TerminatedAt *time.Time
	Error        *string
}
