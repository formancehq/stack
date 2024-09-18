package workflow

import (
	"github.com/formancehq/payments/internal/connectors/engine/activities"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

type UninstallConnector struct {
	ConnectorID models.ConnectorID
}

func (w Workflow) runUninstallConnector(
	ctx workflow.Context,
	uninstallConnector UninstallConnector,
) error {
	if err := workflow.ExecuteChildWorkflow(
		workflow.WithChildOptions(
			ctx,
			workflow.ChildWorkflowOptions{
				TaskQueue:         uninstallConnector.ConnectorID.String(),
				ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
				SearchAttributes: map[string]interface{}{
					SearchAttributeStack: w.stack,
				},
			},
		),
		RunTerminateSchedules,
		TerminateSchedules{
			ConnectorID: uninstallConnector.ConnectorID,
		},
	).Get(ctx, nil); err != nil {
		return errors.Wrap(err, "terminate schedules")
	}

	wg := workflow.NewWaitGroup(ctx)
	errChan := make(chan error, 16)

	wg.Add(1)
	workflow.Go(ctx, func(ctx workflow.Context) {
		defer wg.Done()
		_, err := activities.PluginUninstallConnector(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
		errChan <- err
	})

	wg.Add(1)
	workflow.Go(ctx, func(ctx workflow.Context) {
		defer wg.Done()
		err := activities.StorageSchedulesDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
		errChan <- err
	})

	wg.Add(1)
	workflow.Go(ctx, func(ctx workflow.Context) {
		defer wg.Done()
		err := activities.StorageInstancesDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
		errChan <- err
	})

	wg.Add(1)
	workflow.Go(ctx, func(ctx workflow.Context) {
		defer wg.Done()
		err := activities.StorageTasksTreeDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
		errChan <- err
	})

	wg.Add(1)
	workflow.Go(ctx, func(ctx workflow.Context) {
		defer wg.Done()
		err := activities.StorageBankAccountsDeleteRelatedAccounts(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
		errChan <- err
	})

	wg.Add(1)
	workflow.Go(ctx, func(ctx workflow.Context) {
		defer wg.Done()
		err := activities.StorageAccountsDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
		errChan <- err
	})

	wg.Add(1)
	workflow.Go(ctx, func(ctx workflow.Context) {
		defer wg.Done()
		err := activities.StoragePaymentsDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
		errChan <- err
	})

	wg.Add(1)
	workflow.Go(ctx, func(ctx workflow.Context) {
		defer wg.Done()
		err := activities.StorageStatesDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
		errChan <- err
	})

	wg.Add(1)
	workflow.Go(ctx, func(ctx workflow.Context) {
		defer wg.Done()
		err := activities.StorageWebhooksConfigsDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
		errChan <- err
	})

	wg.Add(1)
	workflow.Go(ctx, func(ctx workflow.Context) {
		defer wg.Done()
		err := activities.StorageWebhooksDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
		errChan <- err
	})

	wg.Wait(ctx)
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	err := activities.StorageConnectorsDelete(infiniteRetryContext(ctx), uninstallConnector.ConnectorID)
	if err != nil {
		return err
	}

	return nil
}

var RunUninstallConnector any

func init() {
	RunUninstallConnector = Workflow{}.runUninstallConnector
}
