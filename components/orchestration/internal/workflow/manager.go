package workflow

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/history/v1"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/sdk/client"
)

var (
	ErrInstanceNotFound = errors.New("Instance not found")
)

const (
	EventSignalName = "event"
)

type Event struct {
	Name string `json:"name"`
}

type Manager struct {
	db             *bun.DB
	temporalClient client.Client
	taskQueue      string
}

func (m *Manager) Create(ctx context.Context, config Config) (*Workflow, error) {

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

func (m *Manager) RunWorkflow(ctx context.Context, id string, variables map[string]string) (Instance, error) {

	workflow := Workflow{}
	if err := m.db.NewSelect().
		Where("id = ?", id).
		Model(&workflow).
		Scan(ctx); err != nil {
		return Instance{}, err
	}

	instance := NewInstance(id)

	if _, err := m.db.
		NewInsert().
		Model(&instance).
		Exec(ctx); err != nil {
		return Instance{}, err
	}

	_, err := m.temporalClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:        instance.ID,
		TaskQueue: m.taskQueue,
	}, Run, Input{
		Instance:  instance,
		Workflow:  workflow,
		Variables: variables,
	})
	if err != nil {
		return Instance{}, err
	}

	return instance, nil
}

func (m *Manager) Wait(ctx context.Context, instanceID string) error {
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

func (m *Manager) ListWorkflows(ctx context.Context) ([]Workflow, error) {
	workflows := make([]Workflow, 0)
	if err := m.db.NewSelect().
		Model(&workflows).
		Scan(ctx); err != nil {
		return nil, err
	}
	return workflows, nil
}

func (m *Manager) ReadWorkflow(ctx context.Context, id string) (Workflow, error) {
	var workflow Workflow
	if err := m.db.NewSelect().
		Model(&workflow).
		Where("id = ?", id).
		Scan(ctx); err != nil {
		return Workflow{}, err
	}
	return workflow, nil
}

func (m *Manager) PostEvent(ctx context.Context, instanceID string, event Event) error {
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

func (m *Manager) AbortRun(ctx context.Context, instanceID string) error {
	instance := Instance{}
	if err := m.db.NewSelect().
		Model(&instance).
		Where("id = ?", instanceID).
		Scan(ctx); err != nil {
		return errors.Wrap(err, "retrieving workflow execution")
	}

	return m.temporalClient.CancelWorkflow(ctx, instanceID, "")
}

func (m *Manager) ListInstances(ctx context.Context, workflowID string, running bool) ([]Instance, error) {
	instances := make([]Instance, 0)
	query := m.db.NewSelect().Model(&instances)
	if workflowID != "" {
		query = query.Where("workflow_id = ?", workflowID)
	}
	if running {
		query = query.Where("terminated = false")
	}
	if err := query.Scan(ctx); err != nil {
		return nil, errors.Wrap(err, "retrieving workflow")
	}
	return instances, nil
}

type StageHistory struct {
	Name         string         `json:"name"`
	Input        map[string]any `json:"input"`
	Error        string         `json:"error,omitempty"`
	Terminated   bool           `json:"terminated"`
	StartedAt    time.Time      `json:"startedAt"`
	TerminatedAt *time.Time     `json:"terminatedAt"`
}

func (m *Manager) ReadInstanceHistory(ctx context.Context, instanceID string) ([]StageHistory, error) {
	historyIterator := m.temporalClient.GetWorkflowHistory(ctx, instanceID, "",
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
	Name         string         `json:"name"`
	Input        map[string]any `json:"input"`
	Output       map[string]any `json:"output,omitempty"`
	Error        string         `json:"error,omitempty"`
	Terminated   bool           `json:"terminated"`
	StartedAt    time.Time      `json:"startedAt"`
	TerminatedAt *time.Time     `json:"terminatedAt"`
}

func (m *Manager) ReadStageHistory(ctx context.Context, instanceID string, stage int) ([]ActivityHistory, error) {
	historyIterator := m.temporalClient.GetWorkflowHistory(ctx, fmt.Sprintf("%s-%d", instanceID, stage), "",
		false, enums.HISTORY_EVENT_FILTER_TYPE_ALL_EVENT)
	ret := make([]ActivityHistory, 0)
	for historyIterator.HasNext() {
		event, err := historyIterator.Next()
		if err != nil {
			return nil, err
		}
		switch event.EventType {
		case enums.EVENT_TYPE_ACTIVITY_TASK_SCHEDULED:
			attributes := event.Attributes.(*history.HistoryEvent_ActivityTaskScheduledEventAttributes).ActivityTaskScheduledEventAttributes

			input := make(map[string]any)
			if err := json.Unmarshal(attributes.Input.Payloads[0].Data, &input); err != nil {
				panic(err)
			}

			activityHistory := ActivityHistory{
				Name: attributes.ActivityType.Name,
				Input: map[string]any{
					attributes.ActivityType.Name: input,
				},
				StartedAt: *event.EventTime,
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
						activityHistory.Output = map[string]any{
							attributes.ActivityType.Name: output,
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
			ret = append(ret, activityHistory)
		}
	}
	return ret, nil
}

func (m *Manager) GetInstance(ctx context.Context, instanceID string) (*Instance, error) {
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

func NewManager(db *bun.DB, temporalClient client.Client, taskQueue string) *Manager {
	return &Manager{
		db:             db,
		temporalClient: temporalClient,
		taskQueue:      taskQueue,
	}
}
