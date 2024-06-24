package workflow

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	enums "go.temporal.io/api/enums/v1"
	history "go.temporal.io/api/history/v1"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
)

var (
	ErrInstanceNotFound = errors.New("Instance not found")
	ErrWorkflowNotFound = errors.New("Workflow not found")
)

const (
	EventSignalName = "event"
)

type Event struct {
	Name string `json:"name"`
}

type WorkflowManager struct {
	db                      *bun.DB
	temporalClient          client.Client
	taskQueue               string
	includeSearchAttributes bool
}

func (m *WorkflowManager) Create(ctx context.Context, config Config) (*Workflow, error) {

	if err := config.Validate(); err != nil {
		return nil, err
	}

	workflow := New(config)

	if _, err := m.db.
		NewInsert().
		Model(&workflow).
		Exec(ctx); err != nil {
		return nil, err
	}

	return &workflow, nil
}

func (m *WorkflowManager) DeleteWorkflow(ctx context.Context, id string) error {

	var workflow Workflow

	res, err := m.db.NewUpdate().Model(&workflow).Where("id = ?", id).Set("deleted_at = ?", time.Now()).Exec(ctx)

	if err != nil {
		return err
	}

	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		return ErrWorkflowNotFound
	}

	return nil
}

func (m *WorkflowManager) RunWorkflow(ctx context.Context, id string, variables map[string]string) (*Instance, error) {

	workflow := Workflow{}
	if err := m.db.NewSelect().
		Where("id = ?", id).
		Model(&workflow).
		Scan(ctx); err != nil {
		return nil, err
	}

	searchAttributes := map[string]any{}
	if m.includeSearchAttributes {
		searchAttributes["OrchestrationWorkflowID"] = workflow.ID
	}

	run, err := m.temporalClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		TaskQueue:        m.taskQueue,
		SearchAttributes: searchAttributes,
	}, Initiate, Input{
		Workflow:  workflow,
		Variables: variables,
	})
	if err != nil {
		return nil, err
	}

	instance := &Instance{}
	if err := run.Get(ctx, instance); err != nil {
		return nil, err
	}

	return instance, nil
}

func (m *WorkflowManager) Wait(ctx context.Context, instanceID string) error {
	if err := m.temporalClient.
		GetWorkflow(ctx, instanceID, "").
		Get(ctx, nil); err != nil {
		if errors.Is(err, &serviceerror.NotFound{}) {
			return ErrInstanceNotFound
		}
		return errors.Unwrap(err)
	}
	return nil
}

func (m *WorkflowManager) ListWorkflows(ctx context.Context, query bunpaginate.OffsetPaginatedQuery[any]) (*bunpaginate.Cursor[Workflow], error) {
	sb := m.db.NewSelect()

	return bunpaginate.UsingOffset[any, Workflow](ctx, sb, query,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.Where("deleted_at IS NULL")
		})
}

func (m *WorkflowManager) ReadWorkflow(ctx context.Context, id string) (Workflow, error) {
	var workflow Workflow
	if err := m.db.NewSelect().
		Model(&workflow).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		return Workflow{}, err
	}
	return workflow, nil
}

func (m *WorkflowManager) PostEvent(ctx context.Context, instanceID string, event Event) error {
	stage := Stage{}
	if err := m.db.NewSelect().
		Model(&stage).
		Where("instance_id = ?", instanceID).
		Limit(1).
		OrderExpr("stage desc").
		Scan(ctx); err != nil {
		return errors.Wrap(err, "retrieving workflow")
	}

	err := m.temporalClient.SignalWorkflow(ctx, stage.TemporalWorkflowID(), "", EventSignalName, event)
	if err != nil {
		return errors.Wrap(err, "sending signal to server")
	}

	return nil
}

func (m *WorkflowManager) AbortRun(ctx context.Context, instanceID string) error {
	instance := Instance{}
	if err := m.db.NewSelect().
		Model(&instance).
		Where("id = ?", instanceID).
		Scan(ctx); err != nil {
		return errors.Wrap(err, "retrieving workflow execution")
	}

	return m.temporalClient.CancelWorkflow(ctx, instanceID, "")
}

func (m *WorkflowManager) ListInstances(ctx context.Context, pagination ListInstancesQuery) (*bunpaginate.Cursor[Instance], error) {
	query := m.db.NewSelect()

	return bunpaginate.UsingOffset[ListInstancesOptions, Instance](ctx, query, bunpaginate.OffsetPaginatedQuery[ListInstancesOptions](pagination),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.
				Join("JOIN workflows ON workflows.id = u.workflow_id").
				Where("workflows.deleted_at IS NULL")

			if pagination.Options.WorkflowID != "" {
				query = query.Where("workflows.id = ?", pagination.Options.WorkflowID)
			}
			if pagination.Options.Running {
				query = query.Where("u.terminated = false")
			}

			return query
		})
}

type StageHistory struct {
	Name         string         `json:"name"`
	Input        map[string]any `json:"input"`
	Error        string         `json:"error,omitempty"`
	Terminated   bool           `json:"terminated"`
	StartedAt    time.Time      `json:"startedAt"`
	TerminatedAt *time.Time     `json:"terminatedAt,omitempty"`
}

func (m *WorkflowManager) ReadInstanceHistory(ctx context.Context, instanceID string) ([]StageHistory, error) {

	historyIterator := m.temporalClient.GetWorkflowHistory(ctx, instanceID+"-main", "",
		false, enums.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)
	ret := make([]StageHistory, 0)
	for historyIterator.HasNext() {
		event, err := historyIterator.Next()
		if err != nil {
			return nil, err
		}
		switch event.EventType {
		case enums.EVENT_TYPE_START_CHILD_WORKFLOW_EXECUTION_INITIATED:
			attributes := event.Attributes.(*history.HistoryEvent_StartChildWorkflowExecutionInitiatedEventAttributes)
			input := make(map[string]any)
			if err := json.Unmarshal(attributes.StartChildWorkflowExecutionInitiatedEventAttributes.Input.Payloads[0].Data, &input); err != nil {
				panic(err)
			}
			stageHistory := StageHistory{
				Name:      attributes.StartChildWorkflowExecutionInitiatedEventAttributes.WorkflowType.Name,
				Input:     input,
				StartedAt: *event.EventTime,
			}

			for historyIterator.HasNext() {
				event, err = historyIterator.Next()
				if err != nil {
					return nil, err
				}
				switch event.EventType {
				case enums.EVENT_TYPE_CHILD_WORKFLOW_EXECUTION_TERMINATED:
				case enums.EVENT_TYPE_CHILD_WORKFLOW_EXECUTION_FAILED:
					attributes := event.Attributes.(*history.HistoryEvent_ChildWorkflowExecutionFailedEventAttributes).
						ChildWorkflowExecutionFailedEventAttributes
					stageHistory.Error = attributes.Failure.Message
				case enums.EVENT_TYPE_CHILD_WORKFLOW_EXECUTION_COMPLETED:
				case enums.EVENT_TYPE_CHILD_WORKFLOW_EXECUTION_TIMED_OUT:
					stageHistory.Error = "timeout"
				case enums.EVENT_TYPE_CHILD_WORKFLOW_EXECUTION_CANCELED:
					stageHistory.Error = "canceled"
				default:
					continue
				}
				stageHistory.TerminatedAt = event.EventTime
				stageHistory.Terminated = true
				break
			}
			ret = append(ret, stageHistory)
		}
	}
	return ret, nil
}

type ActivityHistory struct {
	Name          string         `json:"name"`
	Input         map[string]any `json:"input"`
	Output        map[string]any `json:"output,omitempty"`
	Error         string         `json:"error,omitempty"`
	Terminated    bool           `json:"terminated"`
	StartedAt     time.Time      `json:"startedAt"`
	TerminatedAt  *time.Time     `json:"terminatedAt,omitempty"`
	LastFailure   string         `json:"lastFailure,omitempty"`
	Attempt       int            `json:"attempt"`
	NextExecution *time.Time     `json:"nextExecution,omitempty"`
}

func (m *WorkflowManager) ReadStageHistory(ctx context.Context, instanceID string, stage int) ([]*ActivityHistory, error) {
	stageID := fmt.Sprintf("%s-%d", instanceID, stage)
	described, err := m.temporalClient.DescribeWorkflowExecution(ctx, stageID, "")
	if err != nil {
		if _, ok := err.(*serviceerror.NotFound); ok {
			return nil, ErrInstanceNotFound
		}
		panic(err)
	}

	historyIterator := m.temporalClient.GetWorkflowHistory(ctx, stageID, "",
		false, enums.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)
	ret := make([]*ActivityHistory, 0)
	for historyIterator.HasNext() {
		event, err := historyIterator.Next()
		if err != nil {
			return nil, err
		}
		switch event.EventType {
		case enums.EVENT_TYPE_ACTIVITY_TASK_SCHEDULED:
			activityTaskScheduledEventAttributes := event.Attributes.(*history.HistoryEvent_ActivityTaskScheduledEventAttributes).ActivityTaskScheduledEventAttributes
			input := make(map[string]any)
			if err := json.Unmarshal(activityTaskScheduledEventAttributes.Input.Payloads[0].Data, &input); err != nil {
				panic(err)
			}

			activityHistory := &ActivityHistory{
				Name: activityTaskScheduledEventAttributes.ActivityType.Name,
				Input: map[string]any{
					activityTaskScheduledEventAttributes.ActivityType.Name: input,
				},
				StartedAt: *event.EventTime,
				Attempt:   1,
			}

			ret = append(ret, activityHistory)

			if len(described.PendingActivities) > 0 &&
				activityTaskScheduledEventAttributes.ActivityId == described.PendingActivities[0].ActivityId {
				pendingActivity := described.PendingActivities[0]
				if pendingActivity.LastFailure != nil {
					activityHistory.LastFailure = pendingActivity.LastFailure.Message
				}
				activityHistory.Attempt = int(pendingActivity.Attempt)
				activityHistory.NextExecution = pendingActivity.ScheduledTime
				return ret, nil
			}

			for historyIterator.HasNext() {
				event, err = historyIterator.Next()
				if err != nil {
					return nil, err
				}
				switch event.EventType {
				case enums.EVENT_TYPE_ACTIVITY_TASK_CANCELED:
					activityHistory.Error = "cancelled"
				case enums.EVENT_TYPE_ACTIVITY_TASK_COMPLETED:
					result := event.Attributes.(*history.HistoryEvent_ActivityTaskCompletedEventAttributes).ActivityTaskCompletedEventAttributes.Result
					if result != nil && len(result.Payloads) > 0 {
						output := make(map[string]any)
						if err := json.Unmarshal(result.Payloads[0].Data, &output); err != nil {
							panic(err)
						}

						// notes(gfyrag): keep compat with format from ledger v1 (since we have moved to ledger v2 api)
						// maybe we should define proper boundaries on activities, independent of the formance sdk
						// to avoid breaking histories
						switch activityTaskScheduledEventAttributes.ActivityType.Name {
						case "CreateTransaction":
							switch tx := output["data"].(type) {
							case map[string]any:
								tx["txid"] = tx["id"]
								output["data"] = []any{tx}
							}
						}

						activityHistory.Output = map[string]any{
							activityTaskScheduledEventAttributes.ActivityType.Name: output,
						}
					}
				case enums.EVENT_TYPE_ACTIVITY_TASK_TIMED_OUT:
					activityHistory.Error = "timeout"
				case enums.EVENT_TYPE_ACTIVITY_TASK_FAILED:
					activityHistory.Error = event.Attributes.(*history.HistoryEvent_ActivityTaskFailedEventAttributes).
						ActivityTaskFailedEventAttributes.Failure.Message
				default:
					continue
				}
				activityHistory.TerminatedAt = event.EventTime
				activityHistory.Terminated = true
				break
			}
		}
	}
	return ret, nil
}

func (m *WorkflowManager) GetInstance(ctx context.Context, instanceID string) (*Instance, error) {
	occurrence := Instance{}
	err := m.db.NewSelect().
		Model(&occurrence).
		Relation("Statuses").
		Where("id = ?", instanceID).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInstanceNotFound
		}
		return nil, err
	}
	return &occurrence, nil
}

func NewManager(db *bun.DB, temporalClient client.Client, taskQueue string, includeSearchAttributes bool) *WorkflowManager {
	return &WorkflowManager{
		db:                      db,
		temporalClient:          temporalClient,
		taskQueue:               taskQueue,
		includeSearchAttributes: includeSearchAttributes,
	}
}

type ListInstancesOptions struct {
	WorkflowID string
	Running    bool
}

type ListInstancesQuery bunpaginate.OffsetPaginatedQuery[ListInstancesOptions]
