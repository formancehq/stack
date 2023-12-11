package api

import (
	"context"

	api "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/formancehq/orchestration/internal/triggers"
	"github.com/formancehq/orchestration/internal/workflow"
)

//go:generate mockgen -source backend.go -destination backend_generated.go -package api . Backend

type Backend interface {
	CreateTrigger(context context.Context, data triggers.TriggerData) (*triggers.Trigger, error)
	AbortRun(ctx context.Context, id string) error
	Create(ctx context.Context, config workflow.Config) (*workflow.Workflow, error)
	DeleteWorkflow(ctx context.Context, id string) error
	ListInstances(ctx context.Context, pagination workflow.ListInstancesQuery) (*api.Cursor[workflow.Instance], error)
	ListTriggers(ctx context.Context, query triggers.ListTriggersQuery) (*api.Cursor[triggers.Trigger], error)
	ListWorkflows(ctx context.Context, query bunpaginate.OffsetPaginatedQuery[any]) (*api.Cursor[workflow.Workflow], error)
	PostEvent(ctx context.Context, id string, event workflow.Event) error
	GetInstance(ctx context.Context, id string) (*workflow.Instance, error)
	ReadInstanceHistory(ctx context.Context, id string) ([]workflow.StageHistory, error)
	ReadStageHistory(ctx context.Context, instanceID string, stage int) ([]*workflow.ActivityHistory, error)
	ReadWorkflow(ctx context.Context, id string) (workflow.Workflow, error)
	RunWorkflow(ctx context.Context, id string, input map[string]string) (*workflow.Instance, error)
	Wait(ctx context.Context, id string) error
	ListTriggersOccurrences(ctx context.Context, query triggers.ListTriggersOccurrencesQuery) (*api.Cursor[triggers.Occurrence], error)
	DeleteTrigger(ctx context.Context, triggerID string) error
	GetTrigger(ctx context.Context, triggerID string) (*triggers.Trigger, error)
}

func NewDefaultBackend(triggersManager *triggers.TriggerManager, workflowManager *workflow.WorkflowManager) Backend {
	return struct {
		*triggers.TriggerManager
		*workflow.WorkflowManager
	}{
		WorkflowManager: workflowManager,
		TriggerManager:  triggersManager,
	}
}
