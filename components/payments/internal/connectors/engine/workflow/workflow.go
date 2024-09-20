package workflow

import (
	"encoding/json"

	temporalworker "github.com/formancehq/go-libs/temporal"
	"github.com/formancehq/payments/internal/connectors/engine/plugins"
	"github.com/formancehq/payments/internal/connectors/engine/webhooks"
	"go.temporal.io/sdk/client"
)

const (
	SearchAttributeWorkflowID = "PaymentWorkflowID"
	SearchAttributeScheduleID = "PaymentScheduleID"
	SearchAttributeStack      = "Stack"
)

type FromPayload struct {
	ID      string          `json:"id"`
	Payload json.RawMessage `json:"payload"`
}

func (f *FromPayload) GetPayload() json.RawMessage {
	if f == nil {
		return nil
	}
	return f.Payload
}

type Workflow struct {
	temporalClient client.Client

	plugins  plugins.Plugins
	webhooks webhooks.Webhooks

	stack string
}

func New(temporalClient client.Client, plugins plugins.Plugins, webhooks webhooks.Webhooks, stack string) Workflow {
	return Workflow{
		temporalClient: temporalClient,
		plugins:        plugins,
		webhooks:       webhooks,
		stack:          stack,
	}
}

func (w Workflow) DefinitionSet() temporalworker.DefinitionSet {
	return temporalworker.NewDefinitionSet().
		Append(temporalworker.Definition{
			Name: "FetchAccounts",
			Func: w.runFetchNextAccounts,
		}).
		Append(temporalworker.Definition{
			Name: "FetchBalances",
			Func: w.runFetchNextBalances,
		}).
		Append(temporalworker.Definition{
			Name: "FetchExternalAccounts",
			Func: w.runFetchNextExternalAccounts,
		}).
		Append(temporalworker.Definition{
			Name: "FetchOthers",
			Func: w.runFetchNextOthers,
		}).
		Append(temporalworker.Definition{
			Name: "FetchPayments",
			Func: w.runFetchNextPayments,
		}).
		Append(temporalworker.Definition{
			Name: "TerminateSchedules",
			Func: w.runTerminateSchedules,
		}).
		Append(temporalworker.Definition{
			Name: "InstallConnector",
			Func: w.runInstallConnector,
		}).
		Append(temporalworker.Definition{
			Name: "UninstallConnector",
			Func: w.runUninstallConnector,
		}).
		Append(temporalworker.Definition{
			Name: "CreateBankAccount",
			Func: w.runCreateBankAccount,
		}).
		Append(temporalworker.Definition{
			Name: "Run",
			Func: w.run,
		}).
		Append(temporalworker.Definition{
			Name: "RunCreateWebhooks",
			Func: w.runCreateWebhooks,
		}).
		Append(temporalworker.Definition{
			Name: "RunHandleWebhooks",
			Func: w.runHandleWebhooks,
		}).
		Append(temporalworker.Definition{
			Name: "RunStoreWebhookTranslation",
			Func: w.runStoreWebhookTranslation,
		}).
		Append(temporalworker.Definition{
			Name: "RunSendEvents",
			Func: w.runSendEvents,
		})
}
