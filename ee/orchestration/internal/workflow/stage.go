package workflow

import (
	"fmt"
	"time"

	"github.com/formancehq/stack/libs/go-libs/pointer"

	"github.com/uptrace/bun"
)

type Stage struct {
	bun.BaseModel `bun:"table:workflow_instance_stage_statuses"`
	Number        int        `json:"stage" bun:"stage,pk"`
	InstanceID    string     `json:"instanceID" bun:"instance_id,pk"`
	TemporalRunID string     `json:"temporalRunID" bun:"temporal_run_id,pk"`
	StartedAt     time.Time  `json:"startedAt" bun:"started_at"`
	TerminatedAt  *time.Time `json:"terminatedAt,omitempty" bun:"terminated_at"`
	Error         *string    `json:"error,omitempty" bun:"error"`
}

func (s *Stage) SetTerminated(err error, date time.Time) {
	s.TerminatedAt = &date
	if err != nil {
		s.Error = pointer.For(err.Error())
	}
}

func (s *Stage) TemporalWorkflowID() string {
	return fmt.Sprintf("%s-%d", s.InstanceID, s.Number)
}

func NewStage(instanceID, temporalRunID string, number int) Stage {
	return Stage{
		BaseModel:     bun.BaseModel{},
		TemporalRunID: temporalRunID,
		Number:        number,
		InstanceID:    instanceID,
		StartedAt:     time.Now(),
	}
}
