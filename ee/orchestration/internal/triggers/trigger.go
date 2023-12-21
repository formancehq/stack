package triggers

import (
	"fmt"
	"time"

	"github.com/formancehq/orchestration/internal/workflow"

	"github.com/formancehq/stack/libs/go-libs/publish"

	"github.com/expr-lang/expr"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

var (
	ErrMissingEvent      = errors.New("missing event")
	ErrMissingWorkflowID = errors.New("missing workflow id")
)

type ExprCompilationError struct {
	err  error
	expr string
}

func (e ExprCompilationError) Unwrap() error {
	return e.err
}

func (e ExprCompilationError) Error() string {
	return fmt.Sprintf("error compiling expression '%s': %s", e.expr, e.err)
}

func (e ExprCompilationError) Is(err error) bool {
	_, ok := err.(ExprCompilationError)
	return ok
}

func IsExprCompilationError(err error) bool {
	return errors.Is(err, ExprCompilationError{})
}

func newExprCompilationError(expr string, err error) ExprCompilationError {
	return ExprCompilationError{
		err:  err,
		expr: expr,
	}
}

type TriggerData struct {
	Event      string            `json:"event" bun:"event,type:varchar"`
	Filter     *string           `json:"filter,omitempty" bun:"filter,type:varchar"`
	WorkflowID string            `json:"workflowID" bun:"workflow_id,type:varchar"`
	Vars       map[string]string `json:"vars,omitempty" bun:"vars,type:jsonb"`
}

func (t TriggerData) Validate() error {
	if t.Event == "" {
		return ErrMissingEvent
	}
	if t.WorkflowID == "" {
		return ErrMissingWorkflowID
	}
	if t.Filter != nil && *t.Filter != "" {
		_, err := expr.Compile(*t.Filter)
		if err != nil {
			return newExprCompilationError(*t.Filter, err)
		}
	}
	for _, e := range t.Vars {
		_, err := expr.Compile(e)
		if err != nil {
			return newExprCompilationError(e, err)
		}
	}
	return nil
}

type Trigger struct {
	bun.BaseModel `bun:"triggers"`
	TriggerData

	ID        string    `json:"id" bun:"id,type:varchar,pk"`
	CreatedAt time.Time `json:"createdAt" bun:"created_at"`
}

func (t Trigger) GetID() string {
	return t.ID
}

func NewTrigger(data TriggerData) (*Trigger, error) {
	return &Trigger{
		TriggerData: data,
		ID:          uuid.NewString(),
		CreatedAt:   time.Now().Round(time.Microsecond).UTC(),
	}, nil
}

type Occurrence struct {
	bun.BaseModel `bun:"triggers_occurrences"`

	EventID            string               `json:"-" bun:"event_id,pk"`
	TriggerID          string               `json:"triggerID" bun:"trigger_id,pk"`
	WorkflowInstanceID string               `json:"workflowInstanceID" bun:"workflow_instance_id"`
	WorkflowInstance   workflow.Instance    `json:"workflowInstance" bun:"rel:belongs-to,join:workflow_instance_id=id"`
	Date               time.Time            `json:"date" bun:"date"`
	Event              publish.EventMessage `json:"event" bun:"event"`
}

func NewTriggerOccurrence(eventID, triggerID, workflowInstanceID string, event publish.EventMessage) Occurrence {
	return Occurrence{
		TriggerID:          triggerID,
		EventID:            eventID,
		WorkflowInstanceID: workflowInstanceID,
		Date:               time.Now().Round(time.Microsecond).UTC(),
		Event:              event,
	}
}
