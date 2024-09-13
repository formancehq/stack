package workflow

import (
	"context"
	"fmt"

	"github.com/formancehq/payments/internal/connectors/engine/activities"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func (w Workflow) run(
	ctx workflow.Context,
	config models.Config,
	connectorID models.ConnectorID,
	fromPayload *FromPayload,
	taskTree []models.TaskTree,
) error {
	var nextWorkflow interface{}
	var request interface{}
	var capability models.Capability
	for _, task := range taskTree {
		switch task.TaskType {
		case models.TASK_FETCH_ACCOUNTS:
			req := FetchNextAccounts{
				Config:      config,
				ConnectorID: connectorID,
				FromPayload: fromPayload,
			}

			nextWorkflow = RunFetchNextAccounts
			request = req
			capability = models.CAPABILITY_FETCH_ACCOUNTS

		case models.TASK_FETCH_EXTERNAL_ACCOUNTS:
			req := FetchNextExternalAccounts{
				Config:      config,
				ConnectorID: connectorID,
				FromPayload: fromPayload,
			}

			nextWorkflow = RunFetchNextExternalAccounts
			request = req
			capability = models.CAPABILITY_FETCH_EXTERNAL_ACCOUNTS

		case models.TASK_FETCH_OTHERS:
			req := FetchNextOthers{
				Config:      config,
				ConnectorID: connectorID,
				Name:        task.Name,
				FromPayload: fromPayload,
			}

			nextWorkflow = RunFetchNextOthers
			request = req
			capability = models.CAPABILITY_FETCH_OTHERS

		case models.TASK_FETCH_PAYMENTS:
			req := FetchNextPayments{
				Config:      config,
				ConnectorID: connectorID,
				FromPayload: fromPayload,
			}

			nextWorkflow = RunFetchNextPayments
			request = req
			capability = models.CAPABILITY_FETCH_PAYMENTS

		case models.TASK_FETCH_BALANCES:
			req := FetchNextBalances{
				Config:      config,
				ConnectorID: connectorID,
				FromPayload: fromPayload,
			}

			nextWorkflow = RunFetchNextBalances
			request = req
			capability = models.CAPABILITY_FETCH_BALANCES

		case models.TASK_CREATE_WEBHOOKS:
			req := CreateWebhooks{
				Config:      config,
				ConnectorID: connectorID,
				FromPayload: fromPayload,
			}

			nextWorkflow = RunCreateWebhooks
			request = req
			capability = models.CAPABILITY_WEBHOOKS

		default:
			return fmt.Errorf("unknown task type: %v", task.TaskType)
		}

		if task.Periodically {
			// Schedule next workflow every polling duration
			// TODO(polo): context
			var scheduleID string
			if fromPayload == nil {
				scheduleID = fmt.Sprintf("%s-%s", connectorID.String(), capability.String())
			} else {
				scheduleID = fmt.Sprintf("%s-%s-%s", connectorID.String(), capability.String(), fromPayload.ID)
			}
			scheduleHandle, err := w.temporalClient.ScheduleClient().Create(context.Background(), client.ScheduleOptions{
				ID: scheduleID,
				Spec: client.ScheduleSpec{
					Intervals: []client.ScheduleIntervalSpec{
						{
							Every: config.PollingPeriod,
						},
					},
				},
				Action: &client.ScheduleWorkflowAction{
					Workflow: nextWorkflow,
					Args: []interface{}{
						request,
						task.NextTasks,
					},
					TaskQueue: connectorID.String(),
					// Search attributes are used to query workflows
					TypedSearchAttributes: temporal.NewSearchAttributes(
						temporal.NewSearchAttributeKeyKeyword(SearchAttributeScheduleID).ValueSet(scheduleID),
						temporal.NewSearchAttributeKeyKeyword(SearchAttributeStack).ValueSet(w.stack),
					),
				},
				Overlap:            enums.SCHEDULE_OVERLAP_POLICY_SKIP,
				TriggerImmediately: true,
				SearchAttributes: map[string]any{
					SearchAttributeScheduleID: scheduleID,
					SearchAttributeStack:      w.stack,
				},
			})
			if err != nil {
				return err
			}

			err = activities.StorageSchedulesStore(
				infiniteRetryContext(ctx),
				models.Schedule{
					ID:          scheduleHandle.GetID(),
					ConnectorID: connectorID,
					CreatedAt:   workflow.Now(ctx).UTC(),
				})
			if err != nil {
				return err
			}
		} else {
			// Run next workflow immediately
			if err := workflow.ExecuteChildWorkflow(
				workflow.WithChildOptions(
					ctx,
					workflow.ChildWorkflowOptions{
						TaskQueue:         connectorID.String(),
						ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
						SearchAttributes: map[string]interface{}{
							SearchAttributeStack: w.stack,
						},
					},
				),
				nextWorkflow,
				request,
				task.NextTasks,
			).GetChildWorkflowExecution().Get(ctx, nil); err != nil {
				return errors.Wrap(err, "running next workflow")
			}
		}
	}
	return nil
}

var Run any

func init() {
	Run = Workflow{}.run
}
