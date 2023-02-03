package workflow

import (
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type Stage struct {
	bun.BaseModel `bun:"table:workflow_instance_stage_statuses"`
	Number        int        `json:"stage" bun:"stage,pk"`
	InstanceID    string     `json:"instanceID" bun:"instance_id,pk"`
	StartedAt     time.Time  `json:"startedAt"`
	TerminatedAt  *time.Time `json:"terminatedAt,omitempty"`
	Error         string     `json:"error"`
}

func (s *Stage) SetTerminated(err error, date time.Time) {
	s.TerminatedAt = &date
	if err != nil {
		s.Error = err.Error()
	}
}

func (s *Stage) TemporalWorkflowID() string {
	return fmt.Sprintf("%s-%d", s.InstanceID, s.Number)
}

func NewStage(instanceID string, number int) Stage {
	return Stage{
		BaseModel:  bun.BaseModel{},
		Number:     number,
		InstanceID: instanceID,
		StartedAt:  time.Now(),
		Error:      "",
	}
}
